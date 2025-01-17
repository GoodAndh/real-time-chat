package websokcet

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		c.Conn.WriteJSON(message)
	}

}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fiber.NewError(500, fmt.Sprintf("error webscoket: %s", err.Error()))
				break
			}
			fiber.NewError(500, err.Error())
			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}
		
		hub.Broadcast <- msg

	}
}
