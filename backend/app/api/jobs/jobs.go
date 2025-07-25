package jobs

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	JobStatusPending  JobStatus = "pending"
	JobStatusRunning  JobStatus = "running"
	JobStatusComplete JobStatus = "complete"
	JobStatusFailed   JobStatus = "failed"
)

type Job struct {
	ID        string
	Status    JobStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type JobRunner struct {
	jobs map[string]*Job
	mu   sync.Mutex
}

func NewJobRunner() *JobRunner {
	return &JobRunner{
		jobs: make(map[string]*Job),
	}
}

func (r *JobRunner) StartJob() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	jobID := uuid.New().String()
	r.jobs[jobID] = &Job{
		ID:        jobID,
		Status:    JobStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	go r.runJob(jobID)

	return jobID
}

func (r *JobRunner) runJob(jobID string) {
	r.mu.Lock()
	r.jobs[jobID].Status = JobStatusRunning
	r.jobs[jobID].UpdatedAt = time.Now()
	r.mu.Unlock()

	// Simulate work
	time.Sleep(5 * time.Second)

	r.mu.Lock()
	r.jobs[jobID].Status = JobStatusComplete
	r.jobs[jobID].UpdatedAt = time.Now()
	r.mu.Unlock()
}

func (r *JobRunner) GetJobStatus(jobID string) (JobStatus, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, ok := r.jobs[jobID]
	if !ok {
		return "", false
	}

	return job.Status, true
}
