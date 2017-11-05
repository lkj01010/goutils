package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type Logfile struct {
	mu        sync.Mutex
	handle    *os.File
	curName   string
	curSize   int64
	pattern   string
	baseName  string
	frequency int64
	maxSize   int64
	maxFiles  int64
}

func Open(baseName string, frequency, maxSize, maxFiles int64) (io.WriteCloser, error) {
	return open("20060102.150405", baseName, frequency, maxSize, maxFiles)
}

func OpenDaily(baseName string) (io.WriteCloser, error) {
	return open("20060102", baseName, 86400, 0, 0)
}

func open(pattern, baseName string, frequency, maxSize, maxFiles int64) (io.WriteCloser, error) {
	if frequency <= 0 && maxSize <= 0 && maxFiles <= 0 {
		return os.OpenFile(baseName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	}
	lf := &Logfile{pattern: pattern, baseName: baseName, frequency: frequency, maxSize: maxSize, maxFiles: maxFiles}
	if err := lf.rotate(); err != nil {
		lf.Close()
		return nil, err
	}
	go lf.logRotator()
	return lf, nil
}

func (lf *Logfile) Close() error {
	lf.mu.Lock()
	defer lf.mu.Unlock()
	return lf.switchFile(nil, "")
}

func (lf *Logfile) Write(p []byte) (int, error) {
	lf.mu.Lock()
	defer lf.mu.Unlock()
	if lf.handle == nil {
		return 0, os.ErrInvalid
	}
	n, err := lf.handle.Write(p)
	lf.curSize += int64(n)
	if lf.maxSize > 0 && lf.curSize > lf.maxSize {
		lf.rotate()
	}
	return n, err
}

func (lf *Logfile) logRotator() {
	if lf.frequency <= 0 {
		return
	}
	nanoFreq := lf.frequency * 1e9
	for {
		now := time.Now().UnixNano()
		next := (now/nanoFreq)*nanoFreq + nanoFreq
		<-time.After(time.Duration(next - now))
		lf.mu.Lock()
		if lf.handle == nil {
			lf.mu.Unlock()
			return
		} else {
			lf.rotate()
			lf.mu.Unlock()
		}
	}
}

func (lf *Logfile) rotate() error {
	newName := fmt.Sprintf("%s.%s", lf.baseName, time.Now().Format(lf.pattern))
	if newName == lf.curName {
		return nil
	}
	handle, err := os.OpenFile(newName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	lf.switchFile(handle, newName)
	os.Remove(lf.baseName)
	os.Symlink(path.Base(newName), lf.baseName)
	go lf.logPurge()
	return nil
}

func (lf *Logfile) switchFile(newHandle *os.File, newName string) error {
	oldHandle := lf.handle
	lf.handle = newHandle
	lf.curName = newName
	if newHandle != nil {
		fi, err := newHandle.Stat()
		if err != nil {
			return err
		}
		lf.curSize = fi.Size()
	} else {
		lf.curSize = 0
	}
	if oldHandle != nil {
		return oldHandle.Close()
	}
	return nil
}

func (lf *Logfile) logPurge() {
	if lf.maxFiles <= 0 {
		return
	}
	files, _ := filepath.Glob(fmt.Sprintf("%s.*", lf.baseName))
	if n := len(files) - int(lf.maxFiles); n > 0 {
		sort.Strings(files)
		for i := 0; i < n; i++ {
			if files[i] != lf.curName {
				os.Remove(files[i])
			}
		}
	}
}
