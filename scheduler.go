package scheduler

import (
	"time"
	"fmt"
)

//Scheduler struct is created through New func.
// It communicates directly with the queue
// and has no more utility outside this package
// that a pure reference to the actual instance.
type Scheduler struct {
	queue	*queue
	quit	chan bool
	push	chan *Job
	pop		chan *Job
	workers	[]*worker
}

//Creates a new Scheduler and workers as much
// as defined in wsize and apply its channel
// buffer from csize. Also returns a instance.
func New(wsize, csize int) *Scheduler {
	exec := make(chan *Job, csize)
	var ws []*worker
	for i:=0;i<wsize;i++ {
		ws = append(ws, newWorker(exec))
	}

	return &Scheduler{
		queue: newQueue(exec),
		quit: make(chan bool),
		push: make(chan *Job),
		pop: make(chan *Job),
		workers: ws,
	}
}

//Append a now job to the queue with the
// properties from the job structure.
func (s *Scheduler) AppendJob(job *Job) *Job {
	if job.Lapse > 0 && job.Time.IsZero() {
		job.Time = time.Now().Add(job.Lapse)
	}
	//Set job state to alive before sending to
	//the actual queue
	job.alive = true
	s.push <- job
	return job
}

//Removes a job from the queue. Returns false if
// not able to delete it(Like when executed at the
// time set and no lapse specified).
func (s *Scheduler) RemoveJob(job *Job) bool {
	if !job.alive {
		return false
	}
	s.pop <- job
	return true
}

//Stops the queue and destroy the workers.
// Once stopped do not start again. Just make
// a new one if needed.
func (s *Scheduler) Stop() {
	s.quit <- true
}

//Init the Scheduler with a ticker at an specified interval.
func (s *Scheduler) Every(lapse time.Duration) *Scheduler {
	tick := time.NewTicker(lapse)
	go func() {
		defer fmt.Println("Completed!")
		for {
			select {
			case tnow := <- tick.C:
				s.queue.advance(tnow)
			case job := <- s.push:
				s.queue.schedule(job)
			case job := <- s.pop:
				s.queue.unschedule(job)
			case <- s.quit:
				tick.Stop()
				//Destroy all the workers once the timer stops
				// and their action conclude.
				for _, worker := range s.workers {
					worker.quit <- true
				}
				return
			}
		}
	}()
	return s
}
