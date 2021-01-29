package tests

import (
	"fmt"
	"github.com/ElaraLang/elara/base"
	"time"
)

func executeTest(code string) {
	io := newTestIO()
	base.LoadStdLib(io)
	start := time.Now()
	base.Execute("TestCase", mockInputChannel(code), false, io)
	totalTime := time.Since(start)

	fmt.Println("Collected Output!")
	fmt.Println(fmt.Sprintf("%s", io.output))
	fmt.Println("Collected Error!")
	fmt.Println(fmt.Sprintf("%s", io.error))
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

type TestIO struct {
	output []string
	error  []string
}

func newTestIO() *TestIO {
	return &TestIO{
		output: []string{},
		error:  []string{},
	}
}

func (t *TestIO) Print(out string) {
	t.output = append(t.output, out)
}

func (t *TestIO) Println(out string) {
	t.output = append(t.output, out)
}

func (t *TestIO) Printf(format string, args ...interface{}) {
	t.output = append(t.output, fmt.Sprintf(format, args...))
}

func (t *TestIO) Error(out string) {
	t.error = append(t.error, out)
}
