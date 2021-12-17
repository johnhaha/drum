package drum

import (
	"log"
	"sync"
	"time"
)

type RunFunc func() error

//return failure count and err
type OnFail func(int, error)

type RunStatus struct {
	Fail     bool
	TryCount int
	FailAt   time.Time
	Done     chan struct{}
}

var jobLog = make(map[string]*RunStatus)

var jobLogMtx sync.RWMutex

func registerJob(name string) *RunStatus {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	status := &RunStatus{Done: make(chan struct{})}
	jobLog[name] = status
	return status
}

func markStartJob(name string) {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	status := jobLog[name]
	status.TryCount++
	log.Printf("🥁 try start job %v for the %v time", name, status.TryCount)
}

func markStartFail(name string, fail OnFail, err error) {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	status := jobLog[name]
	status.FailAt = time.Now()
	status.Fail = true
	fail(status.TryCount, err)
}

func getRetryTime(name string) time.Duration {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	status := jobLog[name]
	tm := status.TryCount * tryStep
	if tm > maxStep {
		tm = maxStep
	}
	return time.Second * time.Duration(tm)
}

func getStatus(name string) *RunStatus {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	if status, ok := jobLog[name]; ok {
		return status
	}
	return nil
}

func remJob(name string) {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	delete(jobLog, name)
}

func closeJob(name string) {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	status := jobLog[name]
	close(status.Done)
}
