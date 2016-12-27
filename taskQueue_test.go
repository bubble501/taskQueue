package taskQueue_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bubble501/taskQueue"
)

type TestJob struct {
	id int
}

func (job TestJob) Execute() {
	total := 0
	for i := 0; i < job.id; i++ {
		total = total + i
	}
	//  fmt.Printf("the total for %d is %d\n", job.id, total)
}

func TestTaskQueueNormal(*testing.T) {
	queue := taskQueue.New(8, 2000)
	queue.Start()
	for i := 0; i < 300000; i++ {
		queue.AddJob(TestJob{i})
	}

	for {
		time.Sleep(time.Second)
		if queue.IsFree() {
			break
		}
	}

	fmt.Printf("Finished")
}

func TestTaskQueueStop(*testing.T) {
	queue := taskQueue.New(8, 2000)
	queue.Start()
	for i := 0; i < 300000; i++ {
		queue.AddJob(TestJob{i})
	}

	queue.Stop()
}
