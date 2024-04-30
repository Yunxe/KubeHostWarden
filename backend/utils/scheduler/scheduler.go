package scheduler

import (
	"sync"

	"github.com/robfig/cron/v3"
)

var (
	cronScheduler *cron.Cron
	mutex         sync.Mutex
)

func init() {
	cronScheduler = cron.New()
	cronScheduler.Start()
}

func AddJob(spec string, job func()) (cron.EntryID, error) {
	mutex.Lock()
	defer mutex.Unlock()
	return cronScheduler.AddFunc(spec, job)
}

func RemoveJob(id cron.EntryID) {
	mutex.Lock()
	defer mutex.Unlock()
	cronScheduler.Remove(id)
}

func Stop() {
	cronScheduler.Stop()
}
