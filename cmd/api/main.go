package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"jobQueue-go/internal/model"
	"jobQueue-go/internal/queue"
	"jobQueue-go/internal/service"
	"jobQueue-go/pkg/db"
)

func main() {
	r := gin.Default()
	db.Init()

	r.POST("/job", func(c *gin.Context) {
		var req struct {
			Type string `json:"type"`
			Data string `json:"data"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		job := model.Job{
			ID:         uuid.New().String(),
			Type:       req.Type,
			Status:     "queued",
			Data:       req.Data,
			Retries:    0,
			MaxRetries: 3,
		}
		service.CreateJob(job)
		queue.Enqueue(job)
		c.JSON(http.StatusOK, job)
	})

	r.GET("/job/:id", func(c *gin.Context) {
		id := c.Param("id")

		job, err := service.GetJobByID(id)

		if err != nil {
			c.JSON(404, gin.H{"error": "job not found"})
			return
		}
		c.JSON(200, job)
	})
	r.Run(":8080")
}
