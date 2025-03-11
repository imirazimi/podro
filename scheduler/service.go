package scheduler

import (
	"context"
	"interview/pkg"
	"time"
)

type Service struct {
	podroClient PodroClient
}

type PodroClient interface {
	UpdateOrdersStatus(ctx context.Context) error
}

func NewService(podroClient PodroClient) Service {
	return Service{
		podroClient: podroClient,
	}
}

func (s Service) UpdateOrdersStatus() {
	pkg.Logger.Info("UpdateOrdersStatus started")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	err := s.podroClient.UpdateOrdersStatus(ctx)
	if err != nil {
		pkg.Logger.Error(err.Error())
		return
	}
	pkg.Logger.Info("UpdateOrdersStatus ended")
}
