package scheduler

import "github.com/robfig/cron"

type Service interface {
	Start()
}

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) Register(spec string, svc Service) error {
	return s.cron.AddFunc(spec, svc.Start)
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}