package main

import (
	"fmt"
	"github.com/ElaraLang/elara/base"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {
	args := os.Args
	//This isn't really a repl, but it will be. for now, it's close enough in that it will print the output of every expression
	var repl = false
	var fileName = "elara.el"
	for i, arg := range args {
		if arg == "--repl" {
			repl = true
		}
		if arg == "--file" {
			fileName = args[i+1]
		}
	}

	loadStdLib()

	_, input := loadFile(fileName)
	start := time.Now()
	_, lexTime, parseTime, execTime := base.Execute(&fileName, string(input), repl)

	totalTime := time.Since(start)

	fmt.Printf("Lexing took %s\nParsing took %s\nExecution took %s\nExecuted in %s.\n", lexTime, parseTime, execTime, totalTime)
}

func loadStdLib() {
	goPath := os.Getenv("GOPATH")
	filePath := path.Join(goPath, "stdlib/")
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		//Maybe it's not present in GOPATH, let's try in working directory
		cur, err := os.Executable()
		if err != nil {
			panic(err)
		}
		filePath = path.Join(filepath.Dir(cur), "stdlib/")
		filepath.Walk(filePath, loadWalkedFile)
	}()
	filepath.Walk(filePath, loadWalkedFile)
}

func loadWalkedFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		panic(err)
	}
	if info.IsDir() {
		return nil
	}
	_, content := loadFile(path)
	base.Execute(&path, string(content), false)
	return nil
}

func loadFile(fileName string) (string, []byte) {

	input, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}
	return fileName, input
}
