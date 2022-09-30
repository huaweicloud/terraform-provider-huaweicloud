package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListConfsRequest struct {

	// 指定查询集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o ListConfsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListConfsRequest struct{}"
	}

	return strings.Join([]string{"ListConfsRequest", string(data)}, " ")
}
