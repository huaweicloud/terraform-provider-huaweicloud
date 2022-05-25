package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowRecordRuleRequest struct {

	// 规则ID
	Id string `json:"id"`
}

func (o ShowRecordRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRecordRuleRequest struct{}"
	}

	return strings.Join([]string{"ShowRecordRuleRequest", string(data)}, " ")
}
