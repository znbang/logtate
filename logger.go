package logtate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var _ io.WriteCloser = (*Logger)(nil)

type Option struct {
	Path      string // default exename.log or logtate.log
	MaxSize   int    // default 5 MB
	MaxBackup int    // default 5 backups
}

type Logger struct {
	opt  Option
	mu   sync.Mutex
	size int64
	file *os.File
}

func New(option Option) *Logger {
	if option.MaxSize == 0 {
		option.MaxSize = 5
	}
	if option.MaxBackup == 0 {
		option.MaxBackup = 5
	}
	if option.Path == "" {
		option.Path = getName()
	}
	return &Logger{opt: option}
}

func getName() string {
	path, err := os.Executable()
	if err != nil {
		return "logtate.log"
	}
	exe := filepath.Base(path)
	ext := filepath.Ext(exe)
	if ext == ".exe" {
		return strings.TrimSuffix(exe, ext) + ".log"
	}
	return exe + ".log"
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	writeSize := len(p)

	if l.file == nil {
		if err = l.open(); err != nil {
			return 0, err
		}
	}

	if l.opt.MaxSize > 0 && l.size+int64(writeSize) > int64(l.opt.MaxSize)*1048576 {
		if err = l.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = l.file.Write(p)
	l.size += int64(n)

	return n, err
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.close()
}

func (l *Logger) open() error {
	if err := os.MkdirAll(filepath.Dir(l.opt.Path), 0644); err != nil {
		return fmt.Errorf("create directories failed: %w", err)
	}

	f, err := os.OpenFile(l.opt.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open file failed: %w", err)
	}
	l.file = f
	l.size = 0

	if info, err := f.Stat(); err != nil {
		return fmt.Errorf("get file size failed: %w", err)
	} else {
		l.size = info.Size()
	}

	return nil
}

func (l *Logger) rotate() error {
	if err := l.close(); err != nil {
		return err
	}

	ext := filepath.Ext(l.opt.Path)
	prefix := strings.TrimSuffix(l.opt.Path, ext)

	for i := l.opt.MaxBackup; i > 1; i-- {
		oldName := fmt.Sprintf("%s.%d%s", prefix, i-1, ext)
		newName := fmt.Sprintf("%s.%d%s", prefix, i, ext)
		if _, err := os.Stat(oldName); err == nil {
			if err := os.Rename(oldName, newName); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(l.opt.Path, prefix+".1"+ext); err != nil {
		return err
	}

	return l.open()
}

func (l *Logger) close() error {
	if l.file == nil {
		return nil
	}

	err := l.file.Close()
	l.file = nil
	return err
}
