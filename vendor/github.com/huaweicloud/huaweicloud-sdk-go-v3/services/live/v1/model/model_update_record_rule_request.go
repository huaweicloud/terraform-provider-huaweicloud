package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateRecordRuleRequest struct {

	// 规则ID，在创建成功规则后返回
	Id string `json:"id"`

	Body *RecordRuleRequest `json:"body,omitempty"`
}

func (o UpdateRecordRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateRecordRuleRequest struct{}"
	}

	return strings.Join([]string{"UpdateRecordRuleRequest", string(data)}, " ")
}
