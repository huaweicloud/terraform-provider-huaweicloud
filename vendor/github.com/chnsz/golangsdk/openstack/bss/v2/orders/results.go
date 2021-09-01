package orders

import (
	"fmt"
	"time"

	"github.com/chnsz/golangsdk"
)

type UnsubscribeResult struct {
	golangsdk.Result
}

type Order struct {
	ErrorCode string   `json:"error_code"`
	ErrorMsg  string   `json:"error_msg"`
	OrderIDs  []string `json:"order_ids"`
}

func (r UnsubscribeResult) Extract() (*Order, error) {
	var response Order
	err := r.ExtractInto(&response)
	return &response, err
}

type GetResult struct {
	golangsdk.Result
}

type OrderStatus struct {
	ErrorCode string    `json:"error_code"`
	ErrorMsg  string    `json:"error_msg"`
	OrderInfo OrderInfo `json:"order_info"`
}

type OrderInfo struct {
	Status int `json:"status"`
}

func (r GetResult) Extract() (*OrderStatus, error) {
	var response OrderStatus
	err := r.ExtractInto(&response)
	return &response, err
}

func WaitForOrderSuccess(client *golangsdk.ServiceClient, secs int, orderID string) error {
	return golangsdk.WaitFor(secs, func() (bool, error) {
		order := new(OrderStatus)
		_, err := client.Get(getURL(client, orderID), &order, nil)
		if err != nil {
			return false, err
		}
		time.Sleep(15 * time.Second)

		if order.OrderInfo.Status == 5 {
			return true, nil
		}
		if order.OrderInfo.Status == 4 {
			err = fmt.Errorf("Order canceled: %s.", orderID)
			return false, err
		}

		return false, nil
	})
}
