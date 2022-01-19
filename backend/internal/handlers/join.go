package handlers

import (
	"log"

	"github.com/Qwiri/Synq/pkg/handler"
	"github.com/gofiber/websocket/v2"
)

var JoinHandler = &handler.Handler{
	Syntax:      "JOIN",
	Description: "Get's fired when a player joins",
	Handler: func(conn *websocket.Conn) {
		log.Println("placeholder??")
	},
}
