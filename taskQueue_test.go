package taskQueue_test

import (
	"fmt"
	"github.com/bubble501/taskQueue"
	"testing"
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

func TestTastQueue(*testing.T) {
	queue := taskQueue.New(8, 2000)
	queue.Start()
	var job taskQueue.Job
	for i := 0; i < 300000; i++ {
		testjob := TestJob{i}
		job = testjob
		queue.AddJob(job)
	}

	queue.Wait()
	fmt.Printf("Finished")
}
