package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteMuteRuleName 要删除的规则名称
type DeleteMuteRuleName struct {

	// 要删除的静默规则的名称
	Name string `json:"name"`
}

func (o DeleteMuteRuleName) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteMuteRuleName struct{}"
	}

	return strings.Join([]string{"DeleteMuteRuleName", string(data)}, " ")
}
