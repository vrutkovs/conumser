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
	Name string `json:"name"  binding:"required"`
}

type travisMessage struct {
	Branch      string     `json:"branch"  binding:"required"`
	BuildNumber string     `json:"number"  binding:"required"`
	CommitSHA   string     `json:"commit"  binding:"required"`
	Status      string     `json:"status_message"  binding:"required"`
	RepoOwner   string     `json:"owner_name"  binding:"required"`
	Repo        travisRepo `json:"repository"  binding:"required"`
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
	err := c.Bind(&f)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	var m travisMessage
	if err := json.Unmarshal([]byte(f.Payload), &m); err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, nil)
	}
	shortCommit := m.CommitSHA[:7]

	repo := fmt.Sprintf("%s/%s@%s", m.RepoOwner, m.Repo.Name, m.Branch)
	build := m.BuildNumber
	result := getStatus(m.Status)

	message := fmt.Sprintf("*Travis* %s\nBuild %s for commit %s %s", repo, build, shortCommit, result)
	e.tgbot.sendMessage(e.room, message)
	c.JSON(http.StatusOK, "Message sent")
}
