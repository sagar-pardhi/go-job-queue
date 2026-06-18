package jobs

import (
	"context"
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

// CreateJob godoc
//
// @Summary Create a new job
// @Description Creates a new background job and pushes it to Redis
// @Tags Jobs
// @Accept json
// @Produce json
// @Param request body jobs.CreateJobRequest true "Create Job Request"
// @Success 201 {object} jobs.JobResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /jobs [post]
func (h *Handler) CreateJob(c *gin.Context) {
	var req CreateJobRequest

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

// GetJob godoc
//
// @Summary Get job
// @Description Returns a job by ID
// @Tags Jobs
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} jobs.Job
// @Failure 404 {object} map[string]interface{}
// @Router /jobs/{id} [get]
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

// ListJobs godoc
//
// @Summary List jobs
// @Description Returns all jobs
// @Tags Jobs
// @Produce json
// @Success 200 {array} jobs.Job
// @Router /jobs [get]
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

// Metrics godoc
//
// @Summary Get queue metrics
// @Description Returns aggregated job statistics
// @Tags Metrics
// @Produce json
// @Success 200 {object} jobs.MetricsResponse
// @Router /metrics [get]
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
