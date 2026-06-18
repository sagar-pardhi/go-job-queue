package jobs

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	repo  *Repository
	redis *redis.Client
}

func NewHandler(repo *Repository, redis *redis.Client) *Handler {
	return &Handler{
		repo:  repo,
		redis: redis,
	}
}

func (h *Handler) CreateJob(c *gin.Context) {
	var req struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	job := &Job{
		ID:         uuid.New().String(),
		Type:       req.Type,
		Payload:    req.Payload,
		Status:     StatusPending,
		Retries:    0,
		MaxRetries: 3,
	}

	if err := h.repo.Create(job); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Printf("Job created: %s", job.ID)

	h.redis.LPush(
		context.Background(),
		"queue:jobs",
		job.ID,
	)
	log.Printf("Job pushed to queue: %s", job.ID)

	c.JSON(201, job)
}

func (h *Handler) GetJob(c *gin.Context) {
	id := c.Param("id")

	job, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{
			"error": "job not found",
		})
		return
	}

	c.JSON(200, job)
}

func (h *Handler) ListJobs(c *gin.Context) {
	jobs, err := h.repo.List()

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, jobs)
}

func (h *Handler) Metrics(c *gin.Context) {
	metrics, err := h.repo.GetMetrics()

	if err != nil {
		c.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		200,
		metrics,
	)
}
