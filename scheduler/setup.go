package scheduler

import "interview/adapter"

func SetupService(scheduler adapter.Scheduler, podroClient PodroClient) Service {
	svc := NewService(podroClient)
	scheduler.Sch.Every(1).Day().At("00:00").Do(svc.UpdateOrdersStatus)
	// for test
	// scheduler.Sch.Every(5).Second().Do(svc.UpdateOrdersStatus)
	return svc
}
