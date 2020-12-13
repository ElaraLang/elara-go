package base

import (
	"fmt"
	"github.com/mholt/archiver"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"
)

func ExecuteFull(fileName string, scriptMode bool) {
	loadStdLib()

	_, input := loadFile(fileName)
	start := time.Now()
	_, lexTime, parseTime, execTime := Execute(&fileName, string(input), scriptMode)

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

	err = archiver.NewZip().Unarchive(zipPath, to)
	if err != nil {
		panic(err)
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
	Execute(&path, string(content), false)
	return nil
}

func loadFile(fileName string) (string, []byte) {

	input, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}
	return fileName, input
}
