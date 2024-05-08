package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShutdownInstanceRecord struct {

	// 实例id
	InstanceId *string `json:"instance_id,omitempty"`

	// 工作流id
	JobId *string `json:"job_id,omitempty"`
}

func (o ShutdownInstanceRecord) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShutdownInstanceRecord struct{}"
	}

	return strings.Join([]string{"ShutdownInstanceRecord", string(data)}, " ")
}
