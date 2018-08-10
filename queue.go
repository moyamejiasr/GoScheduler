package scheduler

import (
	"time"
	"container/heap"
)

type Function func(...interface{})

type Parameter interface{}

type Parameters []interface{}

//A Job is a list of rules for the queue.
// Job will be executed once and deleted if
// lapse not set. Also time with lapse will be
// seen by the queue as an action that must happen
// at an certain time and after that an specific
// interval of time.
type Job struct {
	index	int
	alive	bool
	Time	time.Time
	Lapse	time.Duration
	Func	Function
	Params	Parameters
}

type jobs []*Job

//An object that contains an ordered list
// of actions to happen and a reference to
// the channel that executes a worker.
type queue struct {
	heap jobs
	exec chan *Job
}

//Creates a new Queue setting it worker channel.
func newQueue(w chan *Job) *queue {
	return &queue{
		exec: w,
	}
}

// Schedule a job for execution at time t.
func (q *queue) schedule(job *Job) {
	heap.Push(&q.heap, job)
}

// Unschedule a certain job.
func (q *queue) unschedule(job *Job) {
	heap.Remove(&q.heap, job.index)
}

// Executes functions previously defined in
// the queue till the time settled.
func (q *queue) advance(t time.Time) {
	for len(q.heap) > 0 && !t.Before(q.heap[0].Time) {
		job := q.heap[0]
		if job.Lapse < 1 {
			heap.Remove(&q.heap, job.index)
		} else {
			job.Time = t.Add(job.Lapse)
			heap.Fix(&q.heap, job.index)
		}
		//Send func to the workers
		q.exec <- job
	}
}

//From here down only heap functions
//__________________________________

func (h jobs) Len() int {
	return len(h)
}

func (h jobs) Less(i, j int) bool {
	return h[i].Time.Before(h[j].Time)
}

func (h jobs) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index, h[j].index = i, j
}

func (h *jobs) Push(x interface{}) {
	data := x.(*Job)
	*h = append(*h, data)
	data.index = len(*h) - 1
}

func (h *jobs) Pop() interface{} {
	n := len(*h)
	data := (*h)[n-1]
	*h = (*h)[:n-1]
	data.index = -1
	return data
}