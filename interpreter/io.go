package interpreter

import (
	"fmt"
	"io"
	"os"
)

type IO interface {
	Print(out string)
	Println(out string)
	Printf(format string, args ...interface{})
	Error(out string)
}

type SysIO struct {
	out io.Writer
	in  io.Reader
	err io.Writer
}

func NewSysIO() *SysIO {
	return &SysIO{
		out: os.Stdout,
		in:  os.Stdin,
		err: os.Stderr,
	}
}

func (s *SysIO) Print(out string) {
	_, err := s.out.Write([]byte(out))
	if err != nil {
		fmt.Println(err)
	}
}

func (s *SysIO) Println(out string) {
	_, err := s.out.Write([]byte(out + "\n"))
	if err != nil {
		fmt.Println(err)
	}
}

func (s *SysIO) Printf(format string, args ...interface{}) {
	out := fmt.Sprintf(format, args...)
	_, err := s.out.Write([]byte(out + "\n"))
	if err != nil {
		fmt.Println(err)
	}
}

func (s *SysIO) Error(out string) {
	_, err := s.err.Write([]byte(out + "\n"))
	if err != nil {
		_ = fmt.Errorf(err.Error())
	}
}
