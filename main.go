package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/TwiN/go-color"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "grownups",
		Short: "A live chat private for everyone",
		Long:  "Dont tread on me",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 || args[0] == "" {
				fmt.Println("Usage: grownups <username>")
				os.Exit(1)
			}

			run(args[0])
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(username string) {
	url := "wss://grownups-server.onrender.com/ws"

	if username == "" {
		fmt.Println("Usage: grownups <username>")
	}

	fmt.Println("Connecting to server... it can take a while :P")

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 120 * time.Second

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Send username to server
	err = conn.WriteMessage(websocket.TextMessage, []byte(username))
	if err != nil {
		fmt.Println("Error sending username:", err)
		return
	}

	fmt.Println(color.InGreen("Connect as " + username))

	scanner := bufio.NewScanner(os.Stdin)

	go readMessages(conn)

	for {
		sendMessage(scanner, conn)
	}
}

func readMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		fmt.Print("\r\033[K")
		fmt.Println(string(msg))
		fmt.Print(color.InGreen("You: "))
	}
}

func sendMessage(scanner *bufio.Scanner, conn *websocket.Conn) {
	fmt.Print(color.InGreen("You: "))
	if scanner.Scan() {
		text := scanner.Text()
		if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}
