package base

import (
	"bufio"
	"fmt"
	"github.com/ElaraLang/elara/interpreter"
	"github.com/mholt/archiver"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"
)

func ExecuteFull(fileName string, scriptMode bool) {
	ioDelegate := interpreter.NewSysIO()
	LoadStdLib(ioDelegate)
	input := loadFile(fileName)
	start := time.Now()
	Execute(fileName, input, scriptMode, ioDelegate)

	totalTime := time.Since(start)

	fmt.Println("===========================")
	fmt.Printf("Executed in %d micro seconds.\n", totalTime.Microseconds())
	fmt.Println("===========================")
}

func LoadStdLib(io interpreter.IO) {
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

	filepath.Walk(elaraPath, loadWalkedFile(io))
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
func loadWalkedFile(io interpreter.IO) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
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
		Execute(path, content, false, io)
		return nil
	}
}

func loadFile(fileName string) chan rune {
	out := make(chan rune)
	go func() {
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanner := bufio.NewReader(file)
		for {
			byte, err := scanner.ReadByte()
			if err != nil {
				out <- -1
				break
			}
			out <- rune(byte)
		}
	}()
	return out
}
