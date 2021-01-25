package base

import (
	"bufio"
	"fmt"
	"github.com/mholt/archiver"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func ExecuteFull(fileName string, scriptMode bool) {
	LoadStdLib()

	input := loadFile(fileName)
	start := time.Now()
	res, lexTime, parseTime, execTime := Execute(&fileName, input, scriptMode)

	totalTime := time.Since(start)

	fmt.Println("===========================")
	fmt.Printf("Lexing took %s\nParsing took %s\nExecution took %s\nExecuted in %s.\n", lexTime, parseTime, execTime, totalTime)
	fmt.Printf("Res: %s\n", res)
	fmt.Println("===========================")
}

func LoadStdLib() {
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

	filepath.Walk(elaraPath, loadWalkedFile)
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
	if filepath.Ext(path) != ".elr" {
		return nil
	}
	content := loadFile(path)
	Execute(&path, content, false)
	return nil
}

func loadFile(fileName string) chan rune {
	out := make(chan rune)
	go func() {
		reader := strings.NewReader(fileName)
		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanRunes)
		for scanner.Scan() {
			out <- rune(scanner.Text()[0])
		}
	}()
	return out
}
