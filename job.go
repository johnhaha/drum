package drum

import (
	"log"
	"sync"
	"time"
)

type RunFunc func() error
type OnFail func(error)

type RunStatus struct {
	Status    bool
	TryCount  int
	SuccessAt time.Time
	FailAt    time.Time
	Done      chan struct{}
}

var jobLog = make(map[string]*RunStatus)

var jobLogMtx sync.RWMutex

func registerJob(name string) *RunStatus {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	status := &RunStatus{}
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

func markStartFail(name string) {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	status := jobLog[name]
	status.FailAt = time.Now()
}

func markJobSuccess(name string) {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	status := jobLog[name]
	status.SuccessAt = time.Now()
	log.Printf("🥁 job %v finish running at %v", name, status.SuccessAt)
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
