package drum

import (
	"context"
	"time"
)

//keep try job until success
func RunJob(ctx context.Context, name string, job RunFunc, fail OnFail) {
	status := registerJob(name)
	markStartJob(name)
	done := make(chan struct{})
	//run job, repeat if err != nil
	err := job()
	if err != nil {
		go func() {
			defer close(done)
			for err != nil {
				markStartFail(name)
				fail(err)
				time.Sleep(getRetryTime(name))
				err = job()
			}
			markJobSuccess(name)
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-status.Done:
				return
			case <-done:
				return
			}
		}
	}
}

//check job status
func CheckJob(name string) *RunStatus {
	status := getStatus(name)
	return status
}

//rem and terminate job
func RemJob(name string) {
	status := getStatus(name)
	if status != nil {
		close(status.Done)
		delete(jobLog, name)
	}
}
