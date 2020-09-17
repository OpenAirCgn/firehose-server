package firehose_server

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logrotation struct {
	BaseFilename     string        // base filename, final name will be baseFilename.YYYYMMDD.csv
	UseDateTree      bool          // use a date based directory structur, i.e YYYY/MM/DD/bla.csv
	Interval         time.Duration // how often to rotate
	currentFileName  string
	currentStartTime time.Time
	currentFile      *os.File
	mu               sync.Mutex
}

func (l *Logrotation) openNewFile() error {
	now := time.Now().UTC()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")
	hour := now.Format("15")
	minute := now.Format("04")
	second := now.Format("05")

	fn := fmt.Sprintf("%s.%s%s%s-%s%s%s.csv", l.BaseFilename, year, month, day, hour, minute, second)

	if l.UseDateTree {
		dir := fmt.Sprintf("%s/%s/%s", year, month, day)
		dir = filepath.FromSlash(dir)
		if err := os.MkdirAll(dir, 0744); err != nil {
			return err
		}
		fn = fmt.Sprintf("%s/%s", dir, fn)
		fn = filepath.FromSlash(fn)
	}
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.currentFile = f
	l.currentFileName = fn
	l.currentStartTime = now
	return nil
}

func (l *Logrotation) Write(bs []byte) (n int, e error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if (time.Since(l.currentStartTime) > l.Interval) || l.currentFile == nil {
		if l.currentFile != nil {
			if err := l.currentFile.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing %s (%v)\n", l.currentFileName, err)
			}
		}
		if e = l.openNewFile(); e != nil {
			return 0, e
		}
	}
	return l.currentFile.Write(bs)
}

func (l *Logrotation) Close() error {
	return l.currentFile.Close()
}
