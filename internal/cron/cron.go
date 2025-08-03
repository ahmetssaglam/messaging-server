package cron

import (
	"fmt"
	"messaging-server/internal/configs"
	log "messaging-server/internal/logging"
	"sync"
	"time"
)

// Cron runs a job every interval, spawning the job in its own goroutine,
// never allowing more than maxConcurrent jobs to overlap.
type Cron struct {
	interval time.Duration
	job      func()
	wg       sync.WaitGroup
	quit     chan struct{}
	sem      chan struct{} // semaphore channel
	running  bool
}

// NewCron returns a Cron that will run job every interval.
func NewCron(job func()) (*Cron, error) {
	// validate the interval
	if configs.AppConfig.CronInterval <= 0 {
		return nil, fmt.Errorf("interval must be > 0, got %d", configs.AppConfig.CronInterval)
	}
	return &Cron{
		interval: time.Duration(configs.AppConfig.CronInterval) * time.Second,
		job:      job,
		sem:      make(chan struct{}, configs.AppConfig.MaxConcurrentJobs),
		running:  false,
	}, nil
}

// Start launches the Cron loop in its own goroutine.
func (c *Cron) Start() {
	// check if the cron is already running
	if c.running {
		log.Logger.Warning("Cron is already running; cannot start again")
		return
	}

	c.quit = make(chan struct{})
	ticker := time.NewTicker(c.interval)

	c.running = true

	go func() {
		for {
			select {
			// wait for the ticker to tick
			case <-ticker.C:
				select {

				// try to acquire a semaphore slot
				// if successful, spawn the job in a goroutine
				case c.sem <- struct{}{}:
					c.wg.Add(1)
					go func() {
						defer func() {
							<-c.sem // release the semaphore
							c.wg.Done()
						}()
						c.job()
					}()

				// if no slots are available, skip this tick
				default:
					log.Logger.Warning("Max concurrency reached; skipping this run")
				}
			// check stop signal
			case <-c.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop signals the Cron to exit and waits for all jobs to complete.
func (c *Cron) Stop() {
	// check if the cron is already stopped
	if !c.running {
		log.Logger.Warning("Cron is not running; cannot stop")
		return
	}

	close(c.quit)
	c.wg.Wait()
	c.running = false
}
