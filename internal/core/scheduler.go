// Package core provides scheduling functionality for the monitoring engine.
package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"dideban/internal/config"

	"github.com/rs/zerolog/log"
)

// ScheduledJob represents a job that can be scheduled for periodic execution.
type ScheduledJob struct {
	// ID is a unique identifier for the job
	ID string

	// Interval is how often the job should run
	Interval time.Duration

	// Task is the function to execute
	Task func(context.Context) error

	// Internal fields
	ticker  *time.Ticker
	cancel  context.CancelFunc
	running bool
}

// Scheduler manages the execution of scheduled monitoring jobs.
// It provides a worker pool for concurrent job execution with proper resource management.
type Scheduler struct {
	config config.SchedulerConfig

	// Job management
	jobs   map[string]*ScheduledJob
	jobsMu sync.RWMutex

	// Worker pool
	workers chan struct{}

	// Lifecycle management
	running bool
	mu      sync.RWMutex
	wg      sync.WaitGroup
	cancel  context.CancelFunc
}

// NewScheduler creates a new scheduler with the given configuration.
//
// Parameters:
//   - cfg: Scheduler configuration
//
// Returns:
//   - *Scheduler: Initialized scheduler instance
func NewScheduler(cfg config.SchedulerConfig) *Scheduler {
	return &Scheduler{
		config:  cfg,
		jobs:    make(map[string]*ScheduledJob),
		workers: make(chan struct{}, cfg.WorkerCount),
	}
}

// Start starts the scheduler and initializes the worker pool.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//
// Returns:
//   - error: Any error that occurred during startup
func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("scheduler is already running")
	}

	// Create cancellable context
	_, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	// Initialize worker pool
	for i := 0; i < s.config.WorkerCount; i++ {
		s.workers <- struct{}{}
	}

	s.running = true
	log.Info().Int("worker_count", s.config.WorkerCount).Msg("Scheduler started")

	return nil
}

// Stop stops the scheduler and all running jobs gracefully.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	log.Info().Msg("Stopping scheduler")

	// Cancel context to stop all jobs
	if s.cancel != nil {
		s.cancel()
	}

	// Stop all jobs
	s.jobsMu.Lock()
	for _, job := range s.jobs {
		s.stopJobUnsafe(job)
	}
	s.jobsMu.Unlock()

	// Wait for all workers to finish
	s.wg.Wait()

	s.running = false
	log.Info().Msg("Scheduler stopped")
}

// AddJob adds a new job to the scheduler and starts it immediately.
//
// Parameters:
//   - job: Job to add and schedule
//
// Returns:
//   - error: Any error that occurred during job addition
func (s *Scheduler) AddJob(job *ScheduledJob) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.running {
		return fmt.Errorf("scheduler is not running")
	}

	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	// Check if job already exists
	if _, exists := s.jobs[job.ID]; exists {
		return fmt.Errorf("job with ID %s already exists", job.ID)
	}

	// Start the job
	if err := s.startJobUnsafe(job); err != nil {
		return fmt.Errorf("failed to start job %s: %w", job.ID, err)
	}

	s.jobs[job.ID] = job
	log.Debug().Str("job_id", job.ID).Dur("interval", job.Interval).Msg("Job added")

	return nil
}

// RemoveJob removes a job from the scheduler and stops it.
//
// Parameters:
//   - jobID: ID of the job to remove
//
// Returns:
//   - error: Any error that occurred during job removal
func (s *Scheduler) RemoveJob(jobID string) error {
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	job, exists := s.jobs[jobID]
	if !exists {
		return fmt.Errorf("job with ID %s not found", jobID)
	}

	s.stopJobUnsafe(job)
	delete(s.jobs, jobID)

	log.Debug().Str("job_id", jobID).Msg("Job removed")
	return nil
}

// GetJobCount returns the number of currently scheduled jobs.
//
// Returns:
//   - int: Number of scheduled jobs
func (s *Scheduler) GetJobCount() int {
	s.jobsMu.RLock()
	defer s.jobsMu.RUnlock()
	return len(s.jobs)
}

// IsRunning returns whether the scheduler is currently running.
//
// Returns:
//   - bool: True if scheduler is running, false otherwise
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// startJobUnsafe starts a job without acquiring locks.
// This method should only be called when appropriate locks are already held.
//
// Parameters:
//   - job: Job to start
//
// Returns:
//   - error: Any error that occurred during job startup
func (s *Scheduler) startJobUnsafe(job *ScheduledJob) error {
	if job.running {
		return fmt.Errorf("job is already running")
	}

	// Create job context
	jobCtx, cancel := context.WithCancel(context.Background())
	job.cancel = cancel

	// Create ticker
	job.ticker = time.NewTicker(job.Interval)

	// Start job goroutine
	s.wg.Add(1)
	go s.runJob(jobCtx, job)

	job.running = true
	return nil
}

// stopJobUnsafe stops a job without acquiring locks.
// This method should only be called when appropriate locks are already held.
//
// Parameters:
//   - job: Job to stop
func (s *Scheduler) stopJobUnsafe(job *ScheduledJob) {
	if !job.running {
		return
	}

	// Cancel job context
	if job.cancel != nil {
		job.cancel()
	}

	// Stop ticker
	if job.ticker != nil {
		job.ticker.Stop()
	}

	job.running = false
}

// runJob runs a single job in its own goroutine.
// It handles the job lifecycle and executes the task at the specified interval.
//
// Parameters:
//   - ctx: Context for cancellation
//   - job: Job to run
func (s *Scheduler) runJob(ctx context.Context, job *ScheduledJob) {
	defer s.wg.Done()
	defer func() {
		if job.ticker != nil {
			job.ticker.Stop()
		}
	}()

	log.Debug().Str("job_id", job.ID).Msg("Job started")

	// Execute immediately on start
	s.executeJobTask(ctx, job)

	// Then execute on ticker
	for {
		select {
		case <-ctx.Done():
			log.Debug().Str("job_id", job.ID).Msg("Job stopped")
			return
		case <-job.ticker.C:
			s.executeJobTask(ctx, job)
		}
	}
}

// executeJobTask executes a job task with proper worker pool management.
// It ensures that the number of concurrent executions doesn't exceed the worker limit.
//
// Parameters:
//   - ctx: Context for cancellation
//   - job: Job whose task to execute
func (s *Scheduler) executeJobTask(ctx context.Context, job *ScheduledJob) {
	// Try to acquire a worker
	select {
	case <-s.workers:
		// Worker acquired, execute task
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			defer func() {
				// Return worker to pool
				s.workers <- struct{}{}
			}()

			// Execute task with retry logic
			s.executeWithRetry(ctx, job)
		}()
	default:
		// No workers available, skip this execution
		log.Warn().Str("job_id", job.ID).Msg("No workers available, skipping job execution")
	}
}

// executeWithRetry executes a job task with retry logic.
//
// Parameters:
//   - ctx: Context for cancellation
//   - job: Job to execute
func (s *Scheduler) executeWithRetry(ctx context.Context, job *ScheduledJob) {
	maxRetries := s.config.MaxRetries

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if ctx.Err() != nil {
			return // Context cancelled
		}

		err := job.Task(ctx)
		if err == nil {
			// Success
			if attempt > 0 {
				log.Info().Str("job_id", job.ID).Int("attempt", attempt+1).Msg("Job succeeded after retry")
			}
			return
		}

		// Log error
		if attempt < maxRetries {
			log.Warn().Str("job_id", job.ID).Int("attempt", attempt+1).Err(err).Msg("Job failed, retrying")

			// Wait before retry with exponential backoff
			backoff := time.Duration(attempt+1) * time.Second
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
				continue
			}
		} else {
			log.Error().Str("job_id", job.ID).Int("attempts", attempt+1).Err(err).Msg("Job failed after all retries")
		}
	}
}
