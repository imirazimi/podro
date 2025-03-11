package podroapp

import (
	"context"
	"encoding/json"
	"fmt"
	"interview/adapter"
	"interview/pkg"
)

type Client struct {
	svc *Service
	otp *adapter.OTP
}

func NewClient(svc *Service, otp *adapter.OTP) *Client {
	return &Client{svc: svc, otp: otp}

}

func (c *Client) UpdateOrdersStatus(ctx context.Context) error {
	_, err := c.svc.UpdateOrdersStatus(ctx, UpdateOrdersStatusReqeust{})
	if err != nil {
		return err
	}
	return nil
}

// main implementaion
// func (c *Client) CallProviderAPI(API string) ([]Order, error) {
// 	resp, err := http.Get(API)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var apiResponse CallProviderAPIResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
// 		return nil, err
// 	}

// 	var orders []Order
// 	for _, data := range apiResponse.Orders {
// 		order := Order{
// 			ID:     data.ID,
// 			Status: OrderStatus(data.Status),
// 		}
// 		orders = append(orders, order)
// 	}

// 	return orders, nil
// }

// note: i use mock call because of incorrect API: need order id
func (c *Client) CallProviderAPI(API string) ([]Order, error) {
	data := `{"data":[{"id":1,"status":"Delivered"},{"id":2,"status":"Inprogress"}]}`
	var apiResponse CallProviderAPIResponse
	if err := json.Unmarshal([]byte(data), &apiResponse); err != nil {
		return []Order{}, err
	}
	var orders []Order
	for _, data := range apiResponse.Orders {
		order := Order{
			ID:     data.ID,
			Status: OrderStatus(data.Status),
		}
		orders = append(orders, order)
	}
	fmt.Println(orders)
	return orders, nil
}

// note: i use mock call for test
func (c *Client) SendSMS(phone string, message string) error {
	pkg.Logger.Info("sms sent")
	return nil
}
