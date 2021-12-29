package drum

import (
	"context"
	"log"
	"time"
)

//keep try job until success
func RunJob(ctx context.Context, name string, job RunFunc, fail OnFail) {
	status := registerJob(name)
	defer remJob(name)

	markStartJob(name)
	//run job, repeat if err != nil
	err := job()
	if err != nil {
		go func() {
			defer closeJob(name)
			for err != nil {
				markStartFail(name, fail, err)
				//wait to retry
				time.Sleep(getRetryTime(name))
				markStartJob(name)
				err = job()
			}
		}()
		for {
			select {
			case <-ctx.Done():
				log.Println("🥁 job done", name)
				return
			case <-status.Done:
				return
			}
		}
	}
	log.Println("🥁 job done", name)
}

//check job status
func CheckJob(name string) *RunStatus {
	status := getStatus(name)
	return status
}

func ConfigStep(config ...DrumConfig) {
	for _, f := range config {
		f()
	}
}
