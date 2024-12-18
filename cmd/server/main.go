package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity (update for production).
	},
}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	// Ensuring connection close when the function exits. preventing resource leak
	defer conn.Close()

	clients[conn] = true

	for {
	
		_,message,err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			delete(clients,conn)
			break
		}
		fmt.Printf("Received: %s\n", message)

		broadcast <- string(message)

		
	}
}

func broadcastMessage(){
	for{
		msg:= <-broadcast

		for client := range clients{
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err!=nil{
				fmt.Println("Error writing message:" ,err)
				client.Close()
				delete(clients,client)
			}
		}
	}
}
func main() {
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("WebSocket server started on :8080")
	go broadcastMessage()
	http.ListenAndServe(":8080", nil)
}
