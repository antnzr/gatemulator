package job

import (
	"sync"
	"time"
)

type JobScheduler interface {
	Schedule(jobId string, delay time.Duration, job func()) error
	Cancel(jobId string) error
}

type JobInfo struct {
	Timer        *time.Timer
	ExecuteCount int
}

type JobManager struct {
	mu   sync.Mutex
	jobs map[string]*JobInfo
}

func NewJobManager() *JobManager {
	return &JobManager{
		jobs: make(map[string]*JobInfo),
	}
}

func (jm *JobManager) ScheduleRepeatable(jobId string, delay time.Duration, fn func()) error {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	jobInfo, exists := jm.jobs[jobId]
	if !exists {
		jobInfo = &JobInfo{
			ExecuteCount: 0,
		}
		jm.jobs[jobId] = jobInfo
	}

	jobInfo.Timer = time.AfterFunc(delay, func() {
		fn()

		jm.mu.Lock()
		defer jm.mu.Unlock()

		jobInfo.ExecuteCount++
	})
	return nil
}

func (jm *JobManager) Delete(jobId string) error {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	if job, ok := jm.jobs[jobId]; ok {
		job.Timer.Stop()
		delete(jm.jobs, jobId)
	}

	return nil
}

func (jm *JobManager) GetJobState(jobId string) int {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	jobInfo, exists := jm.jobs[jobId]
	if !exists {
		return 0
	}

	return jobInfo.ExecuteCount
}
