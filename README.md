# taskQueue

taskQueue is an task queue for golang. it steal the idea even the code from blog [Handling 1 Million Requests per Minute with Go](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/).


The usage of taskQueue is as simple as four steps:

1. Create your own job struct which implement Execute methods.
1. Create the queue.
2. Start the queue.
3. Add job to the queue.
4. Waiting for worker to be finished in the main method if it will exit (optional).

### Example
```golang

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

```
