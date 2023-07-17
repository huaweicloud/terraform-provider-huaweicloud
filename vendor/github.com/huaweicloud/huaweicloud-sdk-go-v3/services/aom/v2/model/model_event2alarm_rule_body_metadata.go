package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Event2alarmRuleBodyMetadata 源数据
type Event2alarmRuleBodyMetadata struct {

	// 用户自定义标签
	CustomField *[]string `json:"customField,omitempty"`
}

func (o Event2alarmRuleBodyMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Event2alarmRuleBodyMetadata struct{}"
	}

	return strings.Join([]string{"Event2alarmRuleBodyMetadata", string(data)}, " ")
}
