package podroapp

import (
	"context"
)

type Repository interface {
	GetProvidersWeaklyReport(ctx context.Context) ([]Report, error)
	UpdateOrders(ctx context.Context, orders []Order) error
	GetOrders(ctx context.Context) ([]Order, error)
	GetProviders(ctx context.Context) ([]Provider, error)
}

type ProviderClient interface {
	CallProviderAPI(API string) ([]Order, error)
}

type SMSClient interface {
	SendSMS(phone string, message string) error
}

type Service struct {
	providerClient ProviderClient
	smsClient      SMSClient
	repo           Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProvidersWeaklyReport(ctx context.Context, req GetProvidersWeaklyReportReqeust) (GetProvidersWeaklyReportResponse, error) {
	reports, err := s.repo.GetProvidersWeaklyReport(ctx)
	if err != nil {
		return GetProvidersWeaklyReportResponse{}, err
	}
	return GetProvidersWeaklyReportResponse{
		WeaklyReports: reports,
	}, nil
}

func (s *Service) UpdateOrdersStatus(ctx context.Context, req UpdateOrdersStatusReqeust) (UpdateOrdersStatusResponse, error) {

	providers, err := s.repo.GetProviders(ctx)
	if err != nil {
		return UpdateOrdersStatusResponse{}, err
	}
	providerOrders := map[uint][]Order{}
	for _, provider := range providers {
		orders, err := s.providerClient.CallProviderAPI(provider.API)
		if err != nil {
			return UpdateOrdersStatusResponse{}, err
		}
		providerOrders[provider.ID] = orders
	}

	orders, err := s.repo.GetOrders(ctx)
	if err != nil {
		return UpdateOrdersStatusResponse{}, err
	}
	for i := range orders {
		for _, providerOrder := range providerOrders[orders[i].ProviderID] {
			if orders[i].ID == providerOrder.ID {
				orders[i].Status = providerOrder.Status
				if providerOrder.Status == PickedUpOrderStatus {
					s.smsClient.SendSMS(orders[i].RecipientPhone, "Your order is pickedup")
				}
				break
			}
		}
	}

	err = s.repo.UpdateOrders(ctx, orders)

	if err != nil {
		return UpdateOrdersStatusResponse{}, err
	}

	return UpdateOrdersStatusResponse{}, nil
}
