package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "strings"
)

// TODO: abstract this so terminal can be created in sdl window too

type terminal struct {
  inputPrompt  string
  outputPrompt string
  inputEnabled bool
}

func (t *terminal) sendMessage(msg string) {
  fmt.Printf("\r%s %s", t.outputPrompt, msg)
  if t.inputEnabled == true {
    fmt.Printf("\r\n%s ", t.inputPrompt)
  }
}

func (t *terminal) readMessage() (string, bool) {
  if t.inputEnabled == false {
    log.Println("Can not read message when input disabled")
    return "", false
  }

  for {
    fmt.Printf("\r%s ", t.inputPrompt)
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    text = strings.TrimSpace(strings.TrimSuffix(text, "\n"))
    if text == "" {
      continue
    }
    return text, true
  }
}
