package orders

import (
	"github.com/huaweicloud/golangsdk"
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
