package taskQueue

import (
	"time"
)

//Job define the method Execute which should be provided by the concrete job struct.
type Job interface {
	Execute()
}

//Worker represent an internal goroutine in the taskQueue. it will populate itself
//to the queue(workerPool) when it starts and then fetch the job from jobChannel.it
//will quit the goroutine once the quit channel is receivable.
type Worker struct {
	workerPool chan chan Job
	jobChannel chan Job
	quit       chan bool
}

func newWorker(workerPool chan chan Job) Worker {
	return Worker{
		workerPool: workerPool,
		jobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w Worker) start() {
	go func() {
		for {
			w.workerPool <- w.jobChannel
			select {
			case job := <-w.jobChannel:
				job.Execute()
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) stop() {
	go func() {
		w.quit <- true
	}()
}

type Queue struct {
	workerPool chan chan Job
	jobQueue   chan Job
	workers    []Worker
}

func New(maxWorkers, maxJobs int) *Queue {
	jobQueue := make(chan Job, maxJobs)
	pool := make(chan chan Job, maxWorkers)
	workers := make([]Worker, maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		workers[i] = newWorker(pool)
	}
	return &Queue{
		workerPool: pool,
		jobQueue:   jobQueue,
		workers:    workers,
	}
}

func (queue *Queue) Start() {
	for _, work := range queue.workers {
		work.start()
	}
	go queue.dispatch()
}

//dispatch will fetch a worker from the workerPool and send the job to that work.
func (queue *Queue) dispatch() {
	for {
		select {
		case job := <-queue.jobQueue:
			go func(job Job) {
				jobChannel := <-queue.workerPool
				jobChannel <- job
			}(job)
		}
	}
}

//Stop will stop all goroutine in the queue.
func (queue *Queue) Stop() {
	for _, work := range queue.workers {
		work.stop()
	}
}

//AddJob will add job to queue's jobQueue.
func (queue *Queue) AddJob(job Job) {
	queue.jobQueue <- job
}

//Wait will wait all job in the queue to be finished. it's main purpose is to avoid
// the exit of whole program before all tasks in queue have been finished.
func (queue *Queue) Wait() {
	for {
		time.Sleep(time.Second)
		if len(queue.jobQueue) == 0 && len(queue.workerPool) == len(queue.workers) {
			break
		}
	}
}
