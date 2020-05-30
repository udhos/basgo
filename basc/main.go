package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const me = "basc"

func main() {

	var baslibModule string
	var baslibImport string

	flag.StringVar(&baslibModule, "baslibModule", "github.com/udhos/basgo@master", "baslib module")
	flag.StringVar(&baslibImport, "baslibImport", "github.com/udhos/basgo/baslib", "baslib package")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("%s: missing input files\n", me)
		flag.Usage()
		os.Exit(1)
	}

	var baseName string
	var catName string

	arg0 := args[0]

	baseName = strings.TrimSuffix(arg0, ".bas")
	if baseName == arg0 {
		log.Fatalf("%s: please input a filename with suffix .bas: %s", me, arg0)
	}
	log.Printf("%s: basename: %s", me, baseName)
	if errMkDir := mkDir(baseName); errMkDir != nil {
		log.Fatalf("%s: %v", me, errMkDir)
	}
	catName = filepath.Join(baseName, arg0)

	errCat := cat(arg0, catName, true)
	if errCat != nil {
		log.Fatalf("%s: cat: %s: %v", me, arg0, errCat)
	}

	for _, arg := range args[1:] {
		log.Printf("%s: arg: %s", me, arg)
		errCat := cat(arg, catName, false)
		if errCat != nil {
			log.Fatalf("%s: cat: %s: %v", me, arg, errCat)
		}
	}

	goOutput := filepath.Join(baseName, baseName+".go")

	if errBasgo := basgo(catName, goOutput, baslibImport); errBasgo != nil {
		log.Fatalf("%s: basgo: %v", me, errBasgo)
	}

	if errBuild := build(baseName, baslibModule); errBuild != nil {
		log.Fatalf("%s: build: %v", me, errBuild)
	}
}

func mkDir(output string) error {
	info, err := os.Stat(output)
	if os.IsNotExist(err) {
		return os.Mkdir(output, 0750)
	}
	if info.Mode().IsDir() {
		return nil
	}
	return fmt.Errorf("could not create directory, file exists: %s", output)
}

func cat(input, output string, create bool) error {
	log.Printf("%s: cat input=%s output=%s", me, input, output)
	fileInput, errInput := os.Open(input)
	if errInput != nil {
		return errInput
	}
	defer fileInput.Close()
	var fileOutput *os.File
	var errOutput error
	if create {
		fileOutput, errOutput = os.Create(output)
	} else {
		fileOutput, errOutput = os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	}
	if errOutput != nil {
		return errOutput
	}
	defer fileOutput.Close()
	_, errCopy := io.Copy(fileOutput, fileInput)
	return errCopy
}

func basgo(input, output, baslibImport string) error {
	log.Printf("%s: basgo: input=%s output=%s baslibImport=%s", me, input, output, baslibImport)
	fileInput, errInput := os.Open(input)
	if errInput != nil {
		return errInput
	}
	defer fileInput.Close()
	fileOutput, errOutput := os.Create(output)
	if errOutput != nil {
		return errOutput
	}
	defer fileOutput.Close()
	cmd := exec.Command("basgo-build", "-baslibImport", baslibImport)
	cmd.Stdin = fileInput
	cmd.Stdout = fileOutput
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func build(dir, baslibModule string) error {
	log.Printf("%s: build: dir=%s baslibModule=%s", me, dir, baslibModule)
	oldDir, errDir := os.Getwd()
	if errDir != nil {
		return errDir
	}
	defer os.Chdir(oldDir)
	if errChdir := os.Chdir(dir); errChdir != nil {
		return errChdir
	}

	_, err := os.Stat("go.mod")
	if os.IsNotExist(err) {
		cmdModInit := exec.Command("go", "mod", "init", dir)
		cmdModInit.Stdin = os.Stdin
		cmdModInit.Stdout = os.Stdout
		cmdModInit.Stderr = os.Stderr
		if errModInit := cmdModInit.Run(); errModInit != nil {
			return errModInit
		}
	} else {
		log.Printf("%s: build: go.mod exists", me)
	}

	cmdGet := exec.Command("go", "get", baslibModule)
	cmdGet.Stdin = os.Stdin
	cmdGet.Stdout = os.Stdout
	cmdGet.Stderr = os.Stderr
	if errGet := cmdGet.Run(); errGet != nil {
		return errGet
	}

	cmdBuild := exec.Command("go", "build")
	cmdBuild.Stdin = os.Stdin
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	if errBuild := cmdBuild.Run(); errBuild != nil {
		return errBuild
	}

	return nil
}
