package main

import (
  "fmt"
  "time"
  "syscall"
  "bufio"
  "os"
  "strings"
)

var (
  logFile = "log"
  updated chan string
  deleted chan string
)

func main() {
  syscall.Mkfifo(logFile, 0600)

  f, err := os.OpenFile(logFile, os.O_CREATE, os.ModeNamedPipe)
  if err != nil {
    panic(err)
  }

  updated = make(chan string)
  deleted = make(chan string)

  reader := bufio.NewReader(f)

  go readChannel()
  readLog(reader)
}


func readLog(reader *bufio.Reader) {
  for {
    line, err := reader.ReadString('\n')
    if err != nil {
      time.Sleep(50 * time.Millisecond)
      continue
    }

    if strings.HasPrefix(line, "updated: ") {
      updated <- strings.TrimPrefix(line, "updated: ")

    } else if strings.HasPrefix(line, "deleted: ") {
      deleted <- strings.TrimPrefix(line, "deleted: ")
    }
  }
}


func readChannel() {
  for {
    select {
    case c := <- updated:
      fmt.Println("update ", c)

    case c := <- deleted:
      fmt.Println("delete ", c)

    case <-time.After(1 * time.Second):
    }
  }
}
