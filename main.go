package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dexterlb/mpvipc"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Please only pass a pipe path.")
		return
	}

	conn := mpvipc.NewConnection(args[1])
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn *mpvipc.Connection) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	start(conn)
}
