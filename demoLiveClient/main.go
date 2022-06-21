package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/winsock/gochat/database"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"os"
	"strings"
)

const (
	ColorReset = "\033[0m"
	ColorRed   = "\033[31m"
	ColorGray  = "\033[37m"
)

var username string
var server string

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	conn, _, err := websocket.Dial(ctx, fmt.Sprintf("http://%s/live?username=%s", server, username), nil)
	if err != nil {
		log.Fatalf("Unable to connect to server %s, %s\n", server, err.Error())
	}
	go client(ctx, conn)

	log.Printf("Connected to server %s with username %s\nType /q to exit.\n", server, username)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "/q" {
			break
		}
		// TODO: Send messages with the input
	}

	// This will close the websocket connection
	cancel()
}

func client(ctx context.Context, conn *websocket.Conn) {
	var message database.Message
	for wsjson.Read(ctx, conn, &message) == nil {
		color := ColorRed
		if message.Sender.Username == username {
			color = ColorGray
		}
		fmt.Printf("%s%s%s: %s\n", color, message.Sender.Username, ColorReset, message.Content)
	}
	_ = conn.Close(websocket.StatusNormalClosure, "Client closing")
}

func init() {
	flag.StringVar(&server, "server", "localhost:8080", "Username to connect to the server with")
	flag.StringVar(&username, "username", "", "Username to connect to the server with")
}
