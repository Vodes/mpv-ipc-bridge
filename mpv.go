package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dexterlb/mpvipc"
)

type commandRequest struct {
	Arguments []interface{} `json:"command"`
	ID        int64         `json:"request_id"`
}

var connection *mpvipc.Connection

func handleStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		if strings.HasPrefix(text, "{") {
			result := &commandRequest{}
			err := json.Unmarshal([]byte(text), &result)
			if err != nil {
				log.Printf("Failed to unmarshal command: %s\n", err)
			}

			for i, val := range result.Arguments {
				parseInt, err := strconv.ParseInt(fmt.Sprint(val), 10, 64)
				if err == nil {
					result.Arguments[i] = parseInt
					continue
				}
			}

			if result.Arguments[0] == "set" {
				err := connection.Set(result.Arguments[1].(string), result.Arguments[2])
				if err != nil {
					log.Printf("Failed to set property '%s': %s\n", result.Arguments[1], err)
				}
			} else {
				_, err := connection.Call(result.Arguments...)
				if err != nil {
					log.Printf("Failed to send command: %s\n", err)
				}
			}
			continue
		}
		log.Println("Please pass a valid json command!")
	}
}

func start(conn *mpvipc.Connection) {
	connection = conn
	events, stopListening := conn.NewEventListener()
	go func() {
		conn.WaitUntilClosed()
		stopListening <- struct{}{}
	}()

	go handleStdin()
	for event := range events {
		marshalled, _ := json.Marshal(event)
		fmt.Println(string(marshalled))
	}
}
