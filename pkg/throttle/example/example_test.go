package example_test

import (
	"fmt"
	"runtime"
	"time"

	"github.com/bahelit/ctrl_plus_revise/pkg/throttle"
)

func doWork(th *throttle.Throttle) {
	// Mark this job as done.
	defer th.Done(nil)

	// Simulate some work.
	time.Sleep(time.Second)
}

func ExampleThrottle() {
	// Create a new throttle variable.
	th := throttle.NewThrottle(3)

	for i := 0; i < 10; i++ {
		// Gatekeeper. Do not let more than 3 goroutines to start.
		th.Do()

		go doWork(th)
	}

	// Result is inconsistent without the sleep.
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("Active number of goroutines = %+v\n", runtime.NumGoroutine())

	// Wait for all the jobs to finish.
	th.Finish()

	// Output:
	// Active number of goroutines = 3
}
