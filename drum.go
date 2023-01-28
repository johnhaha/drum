package drum

import (
	"context"
	"log"
	"time"
)

// keep try job until success
func RunJob(ctx context.Context, name string, job RunFunc, fail OnFail, jobSetting ...JobSetting) {
	status, _ := registerJob(name)
	for _, s := range jobSetting {
		s(status)
	}

	defer remJob(name)
	markStartJob(name)
	//run job, repeat if err != nil
	err := job()
	if err != nil {
		go func() {
			lastTry := false
			defer closeJob(name)
			for err != nil {
				markStartFail(name, fail, lastTry, err)
				//wait to retry
				time.Sleep(getRetryTime(name))
				lastTry = markStartJob(name)
				err = job()
				if err != nil && lastTry {
					log.Println("🥁  failed on max time, job quit")
					return
				}
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

func LockRun(name string, job RunFunc) error {
	_, exist := registerJob(name)
	if exist {
		log.Println("🥁 duplicate run job", name)
		return ErrorDuplicateRun{}
	}
	defer remJob(name)
	return job()
}
