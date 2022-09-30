package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowSinkTaskDetailRespTopicsInfo struct {

	// topic名称。
	Topic *string `json:"topic,omitempty"`

	// 分区列表。
	Partitions *[]ShowSinkTaskDetailRespPartitions `json:"partitions,omitempty"`
}

func (o ShowSinkTaskDetailRespTopicsInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowSinkTaskDetailRespTopicsInfo struct{}"
	}

	return strings.Join([]string{"ShowSinkTaskDetailRespTopicsInfo", string(data)}, " ")
}
