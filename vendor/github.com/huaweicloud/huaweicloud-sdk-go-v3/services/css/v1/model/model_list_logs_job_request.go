package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListLogsJobRequest struct {

	// 指定查询集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o ListLogsJobRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListLogsJobRequest struct{}"
	}

	return strings.Join([]string{"ListLogsJobRequest", string(data)}, " ")
}
