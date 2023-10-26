package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowElbDetailRequest Request Object
type ShowElbDetailRequest struct {

	// 指定待查询的集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o ShowElbDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowElbDetailRequest struct{}"
	}

	return strings.Join([]string{"ShowElbDetailRequest", string(data)}, " ")
}
