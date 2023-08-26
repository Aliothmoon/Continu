package build

import (
	"log"
	"testing"
)

func TestSlice(t *testing.T) {
	s := []string{"1"}
	log.Println(s[0], s[1:])
}
func TestLogWriteCloser(t *testing.T) {
	l := NewLogWriteCloser(1)
	l.Write([]byte("ABC\nABC\nAB"))
	l.Write([]byte("ABC\nABC\nAB\n"))
	l.Write([]byte("ABC\n\nAB\n"))
	l.Write([]byte("ABC"))
	l.Write([]byte("ABC"))
	l.Write([]byte("\nABC"))
	l.Close()
}
