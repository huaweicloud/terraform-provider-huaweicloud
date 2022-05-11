package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type InstanceRestartRequsetBody struct {

	// 在线调试时必填。
	Restart *interface{} `json:"restart,omitempty"`
}

func (o InstanceRestartRequsetBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InstanceRestartRequsetBody struct{}"
	}

	return strings.Join([]string{"InstanceRestartRequsetBody", string(data)}, " ")
}
