package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TwiN/go-color"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var  url = "grownups-server.onrender.com"
// var url = "localhost:8080"

var rootCmd = &cobra.Command{
	Use:   "grownups",
	Short: "A live chat private for everyone",
	Long:  "Dont tread on me",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := cmd.Flags().GetBool("list")

		if err != nil {
			log.Fatal(err)
		}

		if list {
			runList()
		} else {
			if len(args) < 1 || args[0] == "" {
				fmt.Println("Usage: grownups <username>")
				os.Exit(1)
			}

			run(args[0])
		}

	},
}

func main() {
	rootCmd.Flags().BoolP("list", "l", false, "List the users logged in")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runList() {
	response, err := http.Get("https://" + url + "/users-count")
	if err != nil {
		log.Fatal("Error making get request: $v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response body: $v", err)
	}

	fmt.Println("Number of active users: " + string(body))
}

func run(username string) {
	if username == "" {
		fmt.Println("Usage: grownups <username>")
	}

	fmt.Println("Connecting to server... it can take a while :P")

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 120 * time.Second

	conn, _, err := dialer.Dial("wss://"+url+"/ws", nil)
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
