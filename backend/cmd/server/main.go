package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Qwiri/Synq/internal/handlers"
	"github.com/Qwiri/Synq/pkg/handler"
	"github.com/Qwiri/Synq/pkg/model"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recov "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)

func init() {
	log.SetHandler(text.Default)
	log.SetLevel(log.DebugLevel)
}

var rooms map[int]*model.Room

var wsRouters = map[string]*handler.Handler{
	"JOIN": handlers.JoinHandler,
}

func onMessage(conn *websocket.Conn, data []byte) (err error) {
	if len(data) == 0 {
		return
	}
	var all bool
	if string(data)[0] == '!' {
		all = true
		data = []byte(string(data)[1:])
	}
	cmd := strings.Split(string(data), " ")[0]
	log.Infof("[ws] got message from client: %s", string(data))

	r, ok := wsRouters[cmd]
	if !ok {
		log.Infof("Could not find a handler for command '%s', continuing with default logic", cmd)
		for client := range clients {
			if client == conn && !all {
				continue
			}
			if err = client.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		}
		return nil
	}

	r.Handler(conn)
	return nil

}

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})
	app.Use(logger.New())
	app.Use(recov.New())

	// SOCKET
	app.Use("/socket", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/socket/:id", websocket.New(func(c *websocket.Conn) {
		lobbyID := c.Params("id")
		log.Infof("[ws] got connection to lobby %s", lobbyID)

		clients[c] = true

		for {
			if _, msg, err := c.ReadMessage(); err != nil {
				log.WithError(err).Warn("[ws] cannot read message from websocket")
				break
			} else if err = onMessage(c, msg); err != nil {
				log.WithError(err).Warn("handling client message resulted in an error")
			}
		}

		delete(clients, c)
	}))

	sc := make(chan os.Signal, 1)
	go func(cancel chan os.Signal) {
		if err := app.Listen(":8080"); err != nil {
			log.WithError(err).Error("Starting WebServer failed")
		}
		sc <- os.Interrupt
	}(sc)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	log.Info("Shutting down WebServer")
	if err := app.Shutdown(); err != nil {
		log.WithError(err).Error("Shutting down failed")
	}
	log.Info("Done!")
}
