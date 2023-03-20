package main

import (
  "log"
)

const SERVER_IP = ":8080"
const PROMPT_IN = "$"
const PROMPT_OUT = ">"

var term terminal

func main() {
  server, err := connectToServer(SERVER_IP)
  if err != nil {
    log.Fatal("Error connecting to server", err)
  }
  if server == nil {
    log.Fatal("somethignw ent wrong")
  }

  term = terminal{
    inputPrompt:  PROMPT_IN,
    outputPrompt: PROMPT_OUT,
    inputEnabled: true,
  }

  for {
    // reader := bufio.NewReader(os.Stdin)
    // fmt.Print(PROMPT_IN, " ")
    // text, _ := reader.ReadString('\n')
    // if strings.TrimSpace(string(text)) == "" {
    // 	continue
    // }
    // // fmt.Fprintf(c, "[USER] "+text+"\n")
    // server.session.sendChatMessage([]string{"all"}, text)
    // if strings.TrimSpace(string(text)) == "STOP" {
    // 	fmt.Println("TCP client exiting...")
    // 	return
    // }

    msg, ok := term.readMessage()
    if ok == false {
      continue
    }
    server.session.sendChatMessage([]string{"all"}, msg)
    if msg == "STOP" {
      log.Println("Client exiting...")
      return
    }
  }
}
