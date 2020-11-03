package main

import (
	"archive/zip"
	"fmt"
	"github.com/ElaraLang/elara/base"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	args := os.Args
	//This isn't really a repl, but it will be. for now, it's close enough in that it will print the output of every expression
	var repl = false
	var fileName *string
	for i, arg := range args {
		if arg == "--repl" {
			repl = true
		}
		if arg == "--file" {
			fileName = &args[i+1]
		}
	}
	if fileName == nil {
		println("Error: no file provided. Please pass a file to execute with the --file argument")
		return
	}

	loadStdLib()

	_, input := loadFile(*fileName)
	start := time.Now()
	_, lexTime, parseTime, execTime := base.Execute(fileName, string(input), repl)

	totalTime := time.Since(start)

	fmt.Printf("Lexing took %s\nParsing took %s\nExecution took %s\nExecuted in %s.\n", lexTime, parseTime, execTime, totalTime)
}

func loadStdLib() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	elaraPath := path.Join(usr.HomeDir, ".elara/")

	err = os.Mkdir(elaraPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	filePath := path.Join(elaraPath, "stdlib/")
	err = os.Mkdir(filePath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	downloadStandardLibrary(filePath)

	filepath.Walk(filePath, loadWalkedFile)
}

func downloadStandardLibrary(to string) {
	zipPath := path.Join(to, "stdlib.zip")
	_, err := os.Stat(zipPath)
	if err == nil && !os.IsNotExist(err) {
		return
	}

	standardLibraryURL := "https://github.com/ElaraLang/elara-stdlib/archive/main.zip"
	resp, err := http.Get(standardLibraryURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(zipPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	unzip(zipPath, to)
}

func unzip(file string, to string) {
	r, err := zip.OpenReader(file)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(to, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(path, filepath.Clean(to)+string(os.PathSeparator)) {
			panic(fmt.Sprintf("%s: illegal file path", path))
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			panic(err)
		}
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		rc, err := f.Open()
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			panic(err)
		}
	}
}

func loadWalkedFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		panic(err)
	}
	if info.IsDir() {
		return nil
	}
	if filepath.Ext(path) != ".el" {
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
