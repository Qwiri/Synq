package handler

import "github.com/gofiber/websocket/v2"

type Handler struct {
	Syntax      string
	Description string
	Handler     func(*websocket.Conn)
}
