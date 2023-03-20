package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type session struct {
	conn           net.Conn
	user           user
	incomingEvents chan event
	outgoingEvents chan event
	closed         bool
}

func (s *session) isUserAssociated() bool {
	return s.user.name != ""
}

func (s *session) associateUser(user user) error {
	s.user = user
	return nil
}

func (s *session) recvIncoming() error {
	buf := make([]byte, 4096)
	n, err := s.conn.Read(buf)
	if err != nil {
		log.Println("Error reading from connection, closing connection")
		s.closed = true
		return nil
	}
	var e event
	err = json.Unmarshal(buf[:n], &e)
	if err != nil {
		log.Println("Error unmarshalling event", err)
		return err
	}
	switch e.EventType {
	case "chat_message":
		var cme chatMessageEvent
		data, ok := e.Data.(map[string]interface{})
		if !ok {
			log.Println("Error decoding chatMessage event data")
			return errors.New("error decoding chatMessageEvent data")
		}
		err := mapstructure.Decode(data, &cme)
		if err != nil {
			log.Println("Error decoding chatMessage event payload", err)
			return err
		}
		if err != nil {
			log.Println("Error decoding chatMessage event payload", err)
		}
		e.Data = &cme

	default:
		log.Fatal("Received event of unknown type")
	}

	s.incomingEvents <- e
	return nil
}

func (s *session) handleEvent() error {
	select {
	case event := <-s.incomingEvents:
		switch e := event.Data.(type) {
		case *chatMessageEvent:
			cme := event.Data.(*chatMessageEvent)
			err := s.handleChatMessage(*cme)
			if err != nil {
				println("Error handling chat message", err)
				return err
			}
		default:
			log.Printf("Received unknown event type: %T", e)
			return errors.New("Received unknown event type")
		}
	}
	return nil
}

func (s *session) sendEvent(e event) {
	s.outgoingEvents <- e
}

func (s *session) handleChatMessage(msg chatMessageEvent) error {
	msg.Message = strings.TrimSpace(strings.TrimSuffix(msg.Message, "\n"))
	log.Printf("[%s]: %s\n", msg.FromUser, msg.Message)
	s.sendChatMessage("Server", "message recieved!")
	return nil
}

func (s *session) sendChatMessage(fromUser string, message string) error {
	var user string
	if s.user.name == "" {
		user = "anonymous"
	} else {
		user = s.user.name
	}
	msg := chatMessageEvent{
		FromUser: fromUser,
		ToUsers:  []string{user},
		Message:  message,
	}
	e := event{
		EventType: "chat_message",
		Data:      msg,
	}
	s.sendEvent(e)
	return nil
}

func (s *session) sendOutgoing() error {
	for e := range s.outgoingEvents {
		eventBytes, err := json.Marshal(e)
		if err != nil {
			log.Printf("Error marshalling event: %v", err)
			return err
		}

		_, err = s.conn.Write(eventBytes)
		if err != nil {
			log.Printf("Error sending event: %v", err)
			return err
		}
	}
	return nil
}
