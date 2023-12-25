package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ModifyInstanceConfigsReq struct {

	// kafka待修改配置列表。
	KafkaConfigs *[]ModifyInstanceConfig `json:"kafka_configs,omitempty"`
}

func (o ModifyInstanceConfigsReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyInstanceConfigsReq struct{}"
	}

	return strings.Join([]string{"ModifyInstanceConfigsReq", string(data)}, " ")
}
