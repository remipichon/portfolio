package main

import "fmt"

type JobAlreadyRunningError struct {
}

func (e *JobAlreadyRunningError) Error() string {
	return fmt.Sprintf("job is already running, wait for completion or attempt to kill it")
}
