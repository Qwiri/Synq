package handlers

import (
	"github.com/gofiber/websocket/v2"
	"github.com/qwiri/parti/pkg/handler"
	"log"
)

var JoinHandler = &handler.Handler{
	Syntax:      "JOIN",
	Description: "Get's fired when a player joins",
	Handler: func(conn *websocket.Conn) {
		log.Println("placeholder??")
	},
}
