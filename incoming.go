package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type genericMessage struct {
	Room    string `json:"room"`
	Message string `json:"message"`
}

func (e *Env) incoming(c *gin.Context) {
	var m genericMessage
	if err := c.BindJSON(&m); err == nil {
		e.tgbot.sendMessage(m.Room, m.Message)
		c.JSON(http.StatusOK, "")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Message validation failed!",
			"error":   err.Error(),
		})
	}
}
