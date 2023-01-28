package drum

import (
	"sync"
	"time"
)

type RunFunc func() error

//return failure count and err

type RunStatus struct {
	Fail     bool
	TryCount int
	FailAt   time.Time
	Done     chan struct{}
	//will max try if * > 0,default is 300
	MaxTry int
	//retry in every try step, default is 5 second
	TryStep int
	//will wait for max time, default is 300 second
	MaxStep int
	//lock run, rem duplicated run
	RunLock bool
}

var jobLog = make(map[string]*RunStatus)

var jobLogMtx sync.RWMutex

func registerJob(name string) (status *RunStatus, exist bool) {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	if s, ok := jobLog[name]; ok {
		return s, true
	}
	status = &RunStatus{Done: make(chan struct{}), TryStep: 5, MaxStep: 300}
	jobLog[name] = status
	return status, false
}

func markStartJob(name string) (lastTry bool) {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	status := jobLog[name]
	status.TryCount++
	lastTry = (status.MaxTry > 0 && status.TryCount >= status.MaxTry)
	return lastTry
}

func markStartFail(name string, fail OnFail, lastTry bool, err error) {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	status := jobLog[name]
	status.FailAt = time.Now()
	status.Fail = true
	fail(RunResult{
		Name:     name,
		TryCount: status.TryCount,
		FailAt:   status.FailAt,
		LastTry:  lastTry,
		Error:    err,
	})
}

func getRetryTime(name string) time.Duration {
	jobLogMtx.Lock()
	defer jobLogMtx.Unlock()
	status := jobLog[name]
	tm := status.TryCount * status.TryStep
	if tm > status.MaxStep {
		tm = status.MaxStep
	}
	return time.Second * time.Duration(tm)
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
