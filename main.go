package main

import (
	"fmt"
	"github.com/winsock/gochat/database"
	"log"
)

func main() {
	fmt.Println("Hello Guild")
	_, err := database.Open()
	if err != nil {
		log.Printf("Error opening database! %s\n", err.Error())
	}
	fmt.Println("Server Started!")
}
