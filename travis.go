package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type travisRepo struct {
	Name string `json:"name"`
}

type travisMessage struct {
	Branch      string     `json:"branch"`
	BuildNumber string     `json:"number"`
	CommitSHA   string     `json:"commit"`
	Status      string     `json:"status_message"`
	RepoOwner   string     `json:"owner_name"`
	Repo        travisRepo `json:"repository"`
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

	var m travisMessage
	if err := c.BindJSON(&m); err == nil {
		shortCommit := m.CommitSHA[:7]

		repo := fmt.Sprintf("%s/%s@%s", m.RepoOwner, m.Repo.Name, m.Branch)
		build := m.BuildNumber
		result := getStatus(m.Status)

		message := fmt.Sprintf("*Travis* %s\nBuild %s for commit %s %s", repo, build, shortCommit, result)
		e.tgbot.sendMessage(e.room, message)
		c.JSON(http.StatusOK, "Message sent")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Message validation failed!",
			"error":   err.Error(),
		})
	}
}
