package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CheckRuleIdListRequestInfo 检查项ID列表
type CheckRuleIdListRequestInfo struct {

	// 检查项ID列表
	CheckRules *[]CheckRuleKeyInfoRequestInfo `json:"check_rules,omitempty"`
}

func (o CheckRuleIdListRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckRuleIdListRequestInfo struct{}"
	}

	return strings.Join([]string{"CheckRuleIdListRequestInfo", string(data)}, " ")
}
