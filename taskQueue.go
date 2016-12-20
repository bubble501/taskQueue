package taskQueue

import (
  "time"
)

type Job interface {
  Execute()
}

type Worker struct {
  workerPool chan chan Job
  jobChannel chan Job
  quit chan bool
}

type Queue struct {
  workerPool chan chan Job
  jobQueue chan Job
  workers []Worker
}


func NewWorker(workerPool chan chan Job) Worker {
  return Worker {
    workerPool: workerPool,
    jobChannel: make(chan Job),
    quit: make(chan bool),
  }
}

func (w Worker) Start() {
  go func () {
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

func (w Worker) Stop() {
  go func() {
    w.quit <- true
  }()
}

func New(maxWorkers, maxJobs int) *Queue {
  jobQueue := make(chan Job, maxJobs)
  pool := make(chan chan Job, maxWorkers)
  workers := make([]Worker, maxWorkers)

  for i:=0; i < maxWorkers; i++ {
    workers[i] = NewWorker(pool)
  }
  return & Queue{
    workerPool: pool,
    jobQueue: jobQueue,
    workers: workers,
  }
}

func (queue *Queue) Start() {
  for _, work := range queue.workers {
    work.Start()
  }

  go queue.dispatch()
}

func (queue *Queue) dispatch() {
  for {
    select {
    case job :=<-queue.jobQueue:
      go func (job Job) {
        jobChannel := <-queue.workerPool
        jobChannel <- job
      }(job)
    }
  }
}

func (queue *Queue) Stop() {
  for _, work := range queue.workers {
    work.Stop()
  }
}

func (queue *Queue) AddJob(job Job) {
  queue.jobQueue <- job
}

func (queue *Queue) Wait() {
  for {
      time.Sleep(time.Second)
      if len(queue.jobQueue) == 0 && len(queue.workerPool) == len(queue.workers) {
        break
      }
  }
}
