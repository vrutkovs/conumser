package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type travisRepo struct {
	Name  string `json:"name"`
	Owner string `json:"owner_name"`
}

type travisMessage struct {
	Branch      string     `json:"branch"`
	BuildNumber string     `json:"number"`
	CommitSHA   string     `json:"commit"`
	Status      string     `json:"status_message"`
	Repo        travisRepo `json:"repository"`
}

type travisForm struct {
	Payload string `form:"payload" binding: "required"`
}

func getStatus(status string) string {
	switch status {
	case "Passed", "Fixed":
		return "passed"
	default:
		return "failed"
	}
}

func (e *Env) travisMessage(c *gin.Context) {
	expectedToken := os.Getenv("WEBHOOK_TOKEN")
	if c.Params.ByName("token") != expectedToken {
		c.JSON(http.StatusBadRequest, "Wrong webhook token")
		return
	}

	var f travisForm
	if err := c.Bind(&f); err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	var m travisMessage
	if err := json.Unmarshal([]byte(f.Payload), &m); err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	shortCommit := m.CommitSHA[:7]

	repo := fmt.Sprintf("%s/%s@%s", m.Repo.Owner, m.Repo.Name, m.Branch)
	build := m.BuildNumber
	result := getStatus(m.Status)

	message := fmt.Sprintf("*Travis* %s\nBuild #%s for commit %s %s", repo, build, shortCommit, result)
	e.tgbot.sendMessage(e.room, message)
	log.Printf(fmt.Sprintf("Posted '%s' to '%s'", message, e.room))
	c.JSON(http.StatusOK, "Message sent")
}
