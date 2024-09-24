package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type InstanceLtsConfigDetailResp struct {

	// 日志类型
	LogType *string `json:"log_type,omitempty"`

	// 日志组id
	LtsGroupId *string `json:"lts_group_id,omitempty"`

	// 日志流id
	LtsStreamId *string `json:"lts_stream_id,omitempty"`

	// 是否开启上传至lts
	Enabled *bool `json:"enabled,omitempty"`
}

func (o InstanceLtsConfigDetailResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InstanceLtsConfigDetailResp struct{}"
	}

	return strings.Join([]string{"InstanceLtsConfigDetailResp", string(data)}, " ")
}
