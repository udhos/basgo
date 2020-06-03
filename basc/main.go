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

	"github.com/udhos/basgo/basgo"
)

const me = "basc"

func main() {
	basgo.ShowVersion(me)

	var baslibModule string
	var baslibImport string
	var getFlags string
	var basgoBuildCommand string
	var output string

	flag.StringVar(&baslibModule, "baslibModule", basgo.DefaultBaslibModule, "baslib module")
	flag.StringVar(&baslibImport, "baslibImport", basgo.DefaultBaslibImport, "baslib package")
	flag.StringVar(&getFlags, "getFlags", "", "go get flags")
	flag.StringVar(&basgoBuildCommand, "basgoBuild", "basgo-build", "basgo-build command")
	flag.StringVar(&output, "output", "", "output file name")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s FILE [flags]\n", me)
		flag.PrintDefaults()
	}

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
	arg0base := filepath.Base(arg0)

	baseName = strings.TrimSuffix(arg0base, ".bas")
	if baseName == arg0base {
		log.Fatalf("%s: please input a filename with suffix .bas: %s", me, arg0)
	}
	log.Printf("%s: basename: %s", me, baseName)
	if errMkDir := mkDir(baseName); errMkDir != nil {
		log.Fatalf("%s: %v", me, errMkDir)
	}
	catName = filepath.Join(baseName, arg0base)

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

	goOutput := filepath.Join(baseName, baseName) + ".go"

	if errBasgo := basgoBuild(basgoBuildCommand, catName, goOutput, baslibImport); errBasgo != nil {
		log.Fatalf("%s: basgo: %v", me, errBasgo)
	}

	if errFmt := gofmt(goOutput); errFmt != nil {
		log.Fatalf("%s: gofmt: %v", me, errFmt)
	}

	if errBuild := buildGo(baseName, baslibModule, strings.Fields(getFlags), output); errBuild != nil {
		log.Fatalf("%s: build: %v", me, errBuild)
	}

	//
	// guess output
	//
	if output == "" {
		output = baseName
	}
	if !filepath.IsAbs(output) {
		output = filepath.Join(baseName, output)
	}
	log.Printf("%s: output: %s", me, output)
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

	if input == output {
		return fmt.Errorf("rejecting same file name: input=%s and output=%s", input, output)
	}
	infoOutput, errStatOutput := os.Stat(output)
	if os.IsExist(errStatOutput) {
		infoInput, _ := os.Stat(input)
		if os.SameFile(infoInput, infoOutput) {
			return fmt.Errorf("rejecting same file: input=%s and output=%s", input, output)
		}
	}

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

func basgoBuild(basgoBuildCommand, input, output, baslibImport string) error {
	log.Printf("%s: basgo-build: command=%s input=%s output=%s baslibImport=%s", me, basgoBuildCommand, input, output, baslibImport)
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
	cmd := exec.Command(basgoBuildCommand, "-baslibImport", baslibImport)
	cmd.Stdin = fileInput
	cmd.Stdout = fileOutput
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gofmt(output string) error {
	log.Printf("%s: gofmt: %s", me, output)
	cmd := exec.Command("gofmt", "-s", "-w", output)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildGo(dir, baslibModule string, getFlags []string, output string) error {
	log.Printf("%s: build: dir=%s baslibModule=%s output=%s", me, dir, baslibModule, output)
	log.Printf("%s: build: entering dir=%s", me, dir)
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

	args := []string{"get"}
	args = append(args, getFlags...)
	args = append(args, baslibModule)

	cmdGet := exec.Command("go", args...)
	cmdGet.Stdin = os.Stdin
	cmdGet.Stdout = os.Stdout
	cmdGet.Stderr = os.Stderr
	if errGet := cmdGet.Run(); errGet != nil {
		return errGet
	}

	args = []string{"build"}

	if output != "" {
		args = append(args, "-o")
		args = append(args, output)
	}

	cmdBuild := exec.Command("go", args...)
	cmdBuild.Stdin = os.Stdin
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	if errBuild := cmdBuild.Run(); errBuild != nil {
		return errBuild
	}

	return nil
}
