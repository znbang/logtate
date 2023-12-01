package logtate

import (
	"fmt"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := New(Option{Path: "hello.test.log"})

	for i := 0; i < 10; i++ {
		fmt.Fprintln(logger, "hello")
		logger.rotate()
	}

	paths := []string{
		"hello.test.log",
		"hello.test.1.log",
		"hello.test.2.log",
		"hello.test.3.log",
		"hello.test.4.log",
		"hello.test.5.log",
	}

	for _, path := range paths {
		_, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			t.Errorf("file not found: %s", path)
		}
	}

	_, err := os.Stat("hello.test.6.log")
	if err == nil {
		t.Errorf("found unexpected file: %s", "hello.test.6.log")
	}

	logger.Close()

	for _, path := range paths {
		os.Remove(path)
	}
}
