package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopLogsRequest Request Object
type StopLogsRequest struct {

	// 指定关闭日志的集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o StopLogsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopLogsRequest struct{}"
	}

	return strings.Join([]string{"StopLogsRequest", string(data)}, " ")
}
