package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// Env holds references to useful objects in router funcs
type Env struct {
	tgbot *tgBot
	room  string
}

func main() {
	// setup webhook listener
	r := gin.Default()

	tgbot, err := createBot()
	if err != nil {
		panic(fmt.Sprintf("Failed to create tgBot: %v", err))
	}
	room := os.Getenv("TELEGRAM_ROOM")
	env := &Env{tgbot: tgbot, room: room}

	// Generic webhook, which would display output in markdown
	r.POST("/incoming/:token", env.incoming)

	// Travis check
	r.POST("/travis/:token", env.travisMessage)

	r.Run(":8080")
}
