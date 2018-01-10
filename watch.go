package main

import (
  "fmt"
  "time"
  "syscall"
  "bufio"
  "os"
  "strings"
)


type Watcher struct {
  added   chan string
  deleted chan string
}

var (
  addedString = "update: added "
  deletedString = "update: removing "
)


func NewWatcher(logFile string) (*Watcher, error) {
  syscall.Mkfifo(logFile, 0600)

  f, err := os.OpenFile(logFile, os.O_CREATE, os.ModeNamedPipe)
  if err != nil {
    return nil, err
  }

  reader := bufio.NewReader(f)

  w := &Watcher{
    added: make(chan string),
    deleted: make(chan string),
  }

  go w.readLog(reader)

  return w, nil
}


func (w *Watcher) readLog(reader *bufio.Reader) {
  for {
    line, err := reader.ReadString('\n')

    if err != nil {
      time.Sleep(100 * time.Millisecond)
      continue
    }

    if strings.HasPrefix(line, addedString) {
      w.added <- strings.TrimPrefix(line, addedString)

    } else if strings.HasPrefix(line, deletedString) {
      w.deleted <- strings.TrimPrefix(line, deletedString)
    }
  }
}


func (w *Watcher) readChannel() {
  for {
    select {
    case c := <- w.added:
      fmt.Println("add", c)

    case c := <- w.deleted:
      fmt.Println("delete", c)

    case <-time.After(100 * time.Millisecond):
    }
  }
}


func main() {
  w, err := NewWatcher("log")

  if err != nil {
    panic(err)
  }

  w.readChannel()
}
