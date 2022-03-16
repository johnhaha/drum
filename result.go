package drum

import (
	"log"
	"time"
)

type RunResult struct {
	Name     string
	TryCount int
	FailAt   time.Time
	LastTry  bool
	Error    error
}

func (result *RunResult) Print() {
	log.Println("fail on", result.Name, "for the", result.TryCount, "time", "err is", result.Error, " last try check", result.LastTry)
}

type OnFail func(RunResult)
