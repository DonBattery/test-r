package model

import "sync"

// Job is an interface of structures with a Run() method
type Job interface {
	Run()
}

// Pool is pool of parallel Jobs
type Pool struct {
	wg   sync.WaitGroup
	jobs []Job
}

// NewPool returns a pointer to a new Pool
func NewPool() *Pool {
	return &Pool{}
}

// AddJob adds one or more Job to the Pool
func (p *Pool) AddJob(j ...Job) {
	p.jobs = append(p.jobs, j...)
}

// InitJobs calls the Run() method of all the Pool's Jobs in parallel.
// For every job runned an unit will bu added to the Pool's WaitGroup.
// When a job is finished, an unit will be distracted from the WaitGroup.
func (p *Pool) InitJobs() {
	for _, job := range p.jobs {
		p.wg.Add(1)
		go func(j Job) {
			defer p.wg.Done()
			j.Run()
		}(job)
	}
}

// WaitAll waits for all the Jobs to finish.
func (p *Pool) WaitAll() {
	p.wg.Wait()
}
