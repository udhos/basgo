package basgo

import (
	"fmt"

	"github.com/udhos/basgo/node"
)

type command interface {
	exec(b *Basgo, printf funcPrintf) (stop bool)
}

func commandNew(n node.Node) (command, error) {
	switch nn := n.(type) {
	case *node.NodeEmpty:
		return &commandEmpty{}, nil
	case *node.NodeEnd:
		return &commandEnd{}, nil
	case *node.NodeList:
		return &commandList{}, nil
	case *node.NodePrint:
		return &commandPrint{expressions: nn.Expressions}, nil
	default:
		return nil, fmt.Errorf("commandNew: unknown command: %v", nn.Name())
	}
}

type commandEmpty struct{}

func (c *commandEmpty) exec(b *Basgo, printf funcPrintf) (stop bool) {
	return
}

type commandEnd struct{}

func (c *commandEnd) exec(b *Basgo, printf funcPrintf) (stop bool) {
	stop = true
	return
}

type commandList struct{}

func (c *commandList) exec(b *Basgo, printf funcPrintf) (stop bool) {
	for _, line := range b.lines {
		printf(line.raw + "\n")
	}

	return
}

type commandPrint struct {
	expressions []node.NodeExp
}

func (c *commandPrint) exec(b *Basgo, printf funcPrintf) (stop bool) {
	for _, e := range c.expressions {
		printf("command.exec: FIXME WRITEME: evaluate exp: %s\n", e.String())
	}
	printf("\n")
	return
}
