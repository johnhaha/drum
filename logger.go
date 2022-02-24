package drum

import "log"

func FailLogger(name string, count int, err error) {
	log.Println("fail on", name, "for the", count, "time", "err is", err)
}
