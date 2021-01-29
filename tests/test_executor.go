package tests

import (
	"fmt"
	"github.com/ElaraLang/elara/base"
	"time"
)

func executeTest(code string) {
	base.LoadStdLib()
	start := time.Now()
	base.Execute("TestCase", mockInputChannel(code), false)
	totalTime := time.Since(start)
	fmt.Println("===========================")
	fmt.Printf("Executed test in %d micro seconds.\n", totalTime.Microseconds())
	fmt.Println("===========================")
}

func mockInputChannel(code string) chan rune {
	channel := make(chan rune, 30)
	go feedToChannel([]rune(code), channel)
	return channel
}

func feedToChannel(code []rune, channel chan rune) {
	for _, v := range code {
		channel <- v
	}
	channel <- rune(-1)
	close(channel)
}
