package controllers

import (
	"net/http"

	"github.com/YL-Tan/GoHomeAi/internal/httpresponse"
	"github.com/YL-Tan/GoHomeAi/internal/workers"
)

type JobController struct {
	Pool *workers.WorkerPool
}

func NewJobController(pool *workers.WorkerPool) *JobController {
	return &JobController{Pool: pool}
}

func (c *JobController) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]int{
		"active_jobs": c.Pool.GetActiveJobs(),
	}

	httpresponse.SendJSON(w, http.StatusOK, true, "Job status fetched successfully", status, nil)
}
