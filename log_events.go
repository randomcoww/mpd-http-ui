//
// get add and remove item events by parsing the mpd log
//

package main

import (
	// "fmt"
	"time"
	"syscall"
	"bufio"
	"os"
	"strings"
)

type LogEvents struct {
	added   chan string
	deleted chan string
}

var (
	addedString = "update: added "
	deletedString = "update: removing "
)

// process to read log to create add and remove events
func NewLogEventParser(logFile string) (*LogEvents, error) {
	syscall.Mkfifo(logFile, 0600)

	f, err := os.OpenFile(logFile, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(f)

	e := &LogEvents{
		added: make(chan string),
		deleted: make(chan string),
	}

	go e.readLog(reader)

	return e, nil
}

// parse logs and send items to add and remove channels
func (e *LogEvents) readLog(reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if strings.Contains(line, addedString) {
			str := strings.Split(line, addedString)
			e.added <- strings.TrimSuffix(str[len(str)-1], "\n")

		} else if strings.Contains(line, deletedString) {
			str := strings.Split(line, deletedString)
			e.deleted <- strings.TrimSuffix(str[len(str)-1], "\n")
		}
	}
}
