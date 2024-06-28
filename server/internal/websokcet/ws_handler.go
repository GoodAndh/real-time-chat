package websokcet

import (
	"realTime/server/utils"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub}
}

func (h *Handler) RegisterRoute(router fiber.Router) {
	router.Post("/ws/createRoom", h.createRoom)
	router.Get("/ws/joinRoom/:roomID", h.roomID, websocket.New(h.wsHandler, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}))
}

type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) createRoom(c *fiber.Ctx) error {
	var req CreateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	return utils.WriteJson(c, 200, "ok", req)

}

func (h *Handler) roomID(c *fiber.Ctx) error {

	roomID := c.Params("roomID")
	clientID := c.Query("userID")
	username := c.Query("username")

	cl := &Client{
		Conn:     &websocket.Conn{},
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	c.Locals("client", cl)

	return c.Next()
}

func (h *Handler) wsHandler(c *websocket.Conn) {
	defer c.Close()

	client, ok := c.Locals("client").(*Client)
	if !ok {
		c.WriteMessage(websocket.CloseMessage,[]byte("invalid client"))
		return
	}

	client.Conn = c

	m := &Message{
		Content:  "New User Just Arrived",
		RoomID:   client.RoomID,
		Username: client.Username,
	}

	h.hub.Register <- client

	h.hub.Broadcast <- m

	go client.writeMessage()

	client.readMessage(h.hub)

}
