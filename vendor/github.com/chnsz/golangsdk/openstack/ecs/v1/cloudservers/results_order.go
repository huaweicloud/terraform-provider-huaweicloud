package cloudservers

import (
	"fmt"
	"time"

	"github.com/chnsz/golangsdk"
)

type OrderResponse struct {
	OrderID   string   `json:"order_id"`
	JobID     string   `json:"job_id"`
	ServerIDs []string `json:"serverIds"`
}

type OrderStatus struct {
	ErrorCode string     `json:"error_code"`
	ErrorMsg  string     `json:"error_msg"`
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Status     int    `json:"status"`
	ResourceId string `json:"resourceId"`
}

type DeleteOrderResponse struct {
	OrderIDs []string `json:"orderIds"`
}

type OrderResult struct {
	golangsdk.Result
}

type DeleteOrderResult struct {
	golangsdk.Result
}

func (r OrderResult) ExtractOrderResponse() (*OrderResponse, error) {
	order := new(OrderResponse)
	err := r.ExtractInto(order)
	return order, err
}

func (r DeleteOrderResult) ExtractDeleteOrderResponse() (*DeleteOrderResponse, error) {
	order := new(DeleteOrderResponse)
	err := r.ExtractInto(order)
	return order, err
}

func (r OrderResult) ExtractOrderStatus() (*OrderStatus, error) {
	order := new(OrderStatus)
	err := r.ExtractInto(order)
	return order, err
}

func WaitForOrderSuccess(client *golangsdk.ServiceClient, secs int, orderID string) error {
	return golangsdk.WaitFor(secs, func() (bool, error) {
		order := new(OrderStatus)
		_, err := client.Get(orderURL(client, orderID), &order, nil)
		if err != nil {
			return false, err
		}
		time.Sleep(5 * time.Second)

		if len(order.Resources) == 0 {
			return false, nil
		}
		instance := order.Resources[0]

		if instance.Status == 1 {
			return true, nil
		}
		if instance.Status == 2 {
			err = fmt.Errorf("Order failed: %s.", orderID)
			return false, err
		}

		return false, nil
	})
}

func WaitForOrderDeleteSuccess(client *golangsdk.ServiceClient, secs int, orderID string) error {
	return golangsdk.WaitFor(secs, func() (bool, error) {
		order := new(OrderStatus)
		_, err := client.Get(orderURL(client, orderID), &order, nil)
		if err != nil {
			return false, err
		}
		time.Sleep(5 * time.Second)

		if len(order.Resources) == 0 {
			return false, nil
		}
		instance := order.Resources[0]

		if instance.Status == 8 {
			return true, nil
		}
		if instance.Status == 2 {
			err = fmt.Errorf("Order failed: %s.", orderID)
			return false, err
		}

		return false, nil
	})
}

func GetOrderResource(client *golangsdk.ServiceClient, orderID string) (interface{}, error) {
	order := new(OrderStatus)
	_, err := client.Get(orderURL(client, orderID), &order, nil)
	if err != nil {
		return false, err
	}
	instance := order.Resources[0]
	if instance.Status == 1 {
		if e := instance.ResourceId; e != "" {
			return e, nil
		}
	}

	return nil, fmt.Errorf("Get Order resource ID error.")
}
