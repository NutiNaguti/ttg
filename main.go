package main

import (
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	phoneNumber := ""

	client, err := NewClient()
	if err != nil {
		log.Error(err)
	}

	events := make(chan UpdateData)
	defer close(events)

	go func(events chan UpdateData) {
		for {
			client.getUpdates(events)
		}
	}(events)

	for {
		select {
		case e := <-events:
			switch e["@type"] {
			case authorizationState:
				switch e[authorizationState] {
				case authorizationStateWaitPhoneNumber:
					go client.AuthorizeViaPhoneNumber(phoneNumber)
				}
			}
		// log.Info(e)
		// time.Sleep(2)
		default:
			time.Sleep(2)
		}
	}
}
