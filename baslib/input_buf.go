package baslib

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

type inputBuf struct {
	reader io.Reader // not shared
	queue  chan []byte

	// unsafe data shared among goroutines
	buf    bytes.Buffer
	broken error

	mutex sync.RWMutex
}

func newInputBuf(r io.Reader) *inputBuf {
	i := &inputBuf{
		reader: r,
		queue:  make(chan []byte),
	}
	go readLoop(i)
	return i
}

// readLoop runs as a goroutine
func readLoop(i *inputBuf) {
	for {
		buf := make([]byte, 10)
		n, err := i.reader.Read(buf)
		if n > 0 {
			i.queue <- buf[:n]
		}
		if err != nil {
			i.mutex.Lock()
			i.broken = err // data shared
			i.mutex.Unlock()
			close(i.queue) // stop servicing channel
			return
		}
		if n < 1 {
			log.Printf("baslib.readLoop: unexpected empty Read()")
			time.Sleep(time.Millisecond * 500)
			continue
		}
	}
}

func (i *inputBuf) getBroken() error {
	i.mutex.RLock()
	err := i.broken
	i.mutex.RUnlock()
	return err
}

func (i *inputBuf) Read(buf []byte) (int, error) {

	for {
		// 1/3. if data in buffer, return it
		i.mutex.Lock()
		if i.buf.Len() > 0 {
			n, err := i.buf.Read(buf)
			i.mutex.Unlock()
			log.Printf("baslib.inputBuf.Read: buffered1: %d", n)
			return n, err
		}
		i.mutex.Unlock()

		// 2/3. if error, return it
		if errBroken := i.getBroken(); errBroken != nil {
			return 0, errBroken
		}

		// 3/3. read more from input stream into buffer
		i.readMore()
	}
}

func (i *inputBuf) ReadBytes(delim byte) (line []byte, err error) {

	for {
		// 1. search delim in current buffer
		i.mutex.Lock()
		buf := i.buf.Bytes()
		index := bytes.IndexByte(buf, delim)
		i.mutex.Unlock()

		log.Printf("baslib.inputBuf.ReadBytes: buf=[%s] index=%d", string(buf), index)

		if index >= 0 {
			// found
			line = make([]byte, index+1)
			_, err = i.Read(line)
			return
		}

		// 2. if error, return it
		if errBroken := i.getBroken(); errBroken != nil {
			if len(buf) > 0 {
				line = make([]byte, len(buf))
				_, err = i.Read(line)
			}
			if err == nil {
				err = errBroken
			}
			return
		}

		// 3. read more from input stream into buffer
		i.readMore()
	}
}

// try grab more data from input stream into empty buffer
func (i *inputBuf) readMore() error {

	// try input stream
	data, ok := <-i.queue
	log.Printf("baslib.inputBuf.readMore: data=%d", len(data))
	if len(data) > 0 {
		// append data into buffer
		i.mutex.Lock()
		_, errWrite := i.buf.Write(data)
		i.mutex.Unlock()
		if errWrite != nil {
			return errWrite
		}
	}

	if errBroken := i.getBroken(); errBroken != nil {
		return errBroken
	}

	if !ok {
		i.mutex.Lock()
		i.broken = fmt.Errorf("baslib.inputBuf.readMore: input channel closed")
		i.mutex.Unlock()
	}

	return nil
}

func (i *inputBuf) Inkey() (byte, bool) {
	i.mutex.Lock()
	b, err := i.buf.ReadByte()
	i.mutex.Unlock()

	if err == nil {
		return b, true // true: found byte
	}

	go i.readMore() // call as goroutine to keep Inkey() non-blocking

	return b, false // false: buffer empty
}
