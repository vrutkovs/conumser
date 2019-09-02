package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type genericMessage struct {
	Room    string `json:"room"`
	Message string `json:"message"`
}

func (e *Env) incoming(c *gin.Context) {
	expectedToken := os.Getenv("WEBHOOK_TOKEN")
	if c.Params.ByName("token") != expectedToken {
		c.JSON(http.StatusBadRequest, "Wrong webhook token")
		return
	}

	var m genericMessage
	if err := c.BindJSON(&m); err == nil {
		e.tgbot.sendMessage(m.Room, m.Message)
		c.JSON(http.StatusOK, "Message sent")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Message validation failed!",
			"error":   err.Error(),
		})
	}
}
