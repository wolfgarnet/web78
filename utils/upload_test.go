package utils

import (
	"testing"
	"os"
)

func TestNewFilename(t *testing.T) {
	p, _ := os.Getwd()
	println(p)
	base := "myfile.jpg"
	file := NewFilename("user-1", base, "session1")
	println(file)
}