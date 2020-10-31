package main

import (
	"elara/base"
	"fmt"
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
	for _, arg := range args {
		if arg == "--repl" {
			repl = true
			break
		}
	}

	loadStdLib()

	fileName, input := loadFile("elara.el")
	start := time.Now()
	_, lexTime, parseTime, execTime := base.Execute(&fileName, string(input), repl)

	totalTime := time.Since(start)

	fmt.Printf("Lexing took %s\nParsing took %s\nExecution took %s\nExecuted in %s.\n", lexTime, parseTime, execTime, totalTime)
}

func loadStdLib() {
	goPath := os.Getenv("GOPATH")
	filePath := path.Join(goPath, "stdlib/")
	filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return nil
		}
		_, content := loadFile(path)
		base.Execute(&path, string(content), false)
		return nil
	})
}

func loadFile(fileName string) (string, []byte) {
	goPath := os.Getenv("GOPATH")
	filePath := path.Join(goPath, fileName)

	input, err := ioutil.ReadFile(filePath)

	if err != nil {
		panic(err)
	}
	return fileName, input
}
