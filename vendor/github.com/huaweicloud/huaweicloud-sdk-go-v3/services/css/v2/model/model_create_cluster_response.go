package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClusterResponse Response Object
type CreateClusterResponse struct {
	Cluster *CreateClusterResp `json:"cluster,omitempty"`

	// 订单号。若创建的是包周期集群，则只返回ordeld参数。
	OrderId        *string `json:"orderId,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateClusterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterResponse struct{}"
	}

	return strings.Join([]string{"CreateClusterResponse", string(data)}, " ")
}
