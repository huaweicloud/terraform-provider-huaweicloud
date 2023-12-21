package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResizeClusterResponse Response Object
type ResizeClusterResponse struct {

	// 任务ID
	JobID *string `json:"jobID,omitempty"`

	// 包周期集群变更规格订单ID
	OrderID        *string `json:"orderID,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ResizeClusterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResizeClusterResponse struct{}"
	}

	return strings.Join([]string{"ResizeClusterResponse", string(data)}, " ")
}
