package websocket

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID uint
}

var Clients = make(map[uint]*Client)
var NotifyChan = make(chan NotificationMessage)
var Lock = sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type NotificationMessage struct {
	UserID  uint   `json:"user_id"`
	Message string `json:"message"`
}

func WebSocketHandler(c *gin.Context) {
	userIDStr := c.Query("user_id")
	id, _ := strconv.ParseUint(userIDStr, 10, 32)
	userID := uint(id)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	Lock.Lock()
	Clients[userID] = &Client{Conn: conn, UserID: userID}
	Lock.Unlock()

	go func() {
		defer conn.Close()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				Lock.Lock()
				delete(Clients, userID)
				Lock.Unlock()
				break
			}
		}
	}()
}




func SendToClient(msg NotificationMessage) {
	Lock.Lock()
	defer Lock.Unlock()

	client, ok := Clients[msg.UserID]
	if !ok {
		fmt.Printf("No WebSocket client connected for UserID=%d\n", msg.UserID)
		return
	}

	fmt.Printf("Sending notification to UserID=%d: %s\n", msg.UserID, msg.Message)

	err := client.Conn.WriteJSON(msg)
	if err != nil {
		fmt.Printf("Failed to send message to UserID=%d: %v\n", msg.UserID, err)
		client.Conn.Close()
		delete(Clients, msg.UserID)
	}
}
