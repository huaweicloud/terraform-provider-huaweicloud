package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListPipelinesRequest struct {

	// 指定查询集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o ListPipelinesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPipelinesRequest struct{}"
	}

	return strings.Join([]string{"ListPipelinesRequest", string(data)}, " ")
}
