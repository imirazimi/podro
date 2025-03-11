package adapter

import (
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	Sch *gocron.Scheduler
}

func NewScheduler() (Scheduler, error) {
	location, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return Scheduler{}, err
	}
	return Scheduler{
		Sch: gocron.NewScheduler(location),
	}, nil
}

func (s Scheduler) Start() {
	s.Sch.StartAsync()
}

func (s Scheduler) ShutDown() {
	s.Sch.Stop()
}
