//
// get add and remove item events by parsing the mpd log
//

package server

import (
	"bufio"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type MpdLogEvents struct {
	AddEvent    chan string
	DeleteEvent chan string
}

const (
	logAddedString   = "update: added "
	logDeletedString = "update: removing "
)

// process to read log to create add and remove events
func NewMpdLogReader(logFile string) (*MpdLogEvents, error) {
	logrus.Infof("Create MPD log pipe: %s", logFile)

	// os.Remove(logFile)
	syscall.Mkfifo(logFile, 0600)
	// err := syscall.Mkfifo(logFile, 0600)
	// if err != nil {
	// 	return nil, err
	// }

	f, err := os.OpenFile(logFile, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		return nil, err
	}

	e := &MpdLogEvents{
		AddEvent:    make(chan string),
		DeleteEvent: make(chan string),
	}

	go e.run(bufio.NewReader(f))
	return e, nil
}

// parse logs and send items to add and remove channels
func (e *MpdLogEvents) run(reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')
		logrus.Infof("%s", line)

		if err != nil {
			logrus.Errorf("Error reading from pipe: %v", err)
			time.Sleep(1000 * time.Millisecond)
		}

		if strings.Contains(line, logAddedString) {
			str := strings.Split(line, logAddedString)
			e.AddEvent <- strings.TrimSuffix(str[len(str)-1], "\n")

		} else if strings.Contains(line, logDeletedString) {
			str := strings.Split(line, logDeletedString)
			e.DeleteEvent <- strings.TrimSuffix(str[len(str)-1], "\n")
		}
	}
}
