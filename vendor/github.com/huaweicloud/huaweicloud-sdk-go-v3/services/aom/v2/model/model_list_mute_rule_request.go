package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListMuteRuleRequest struct {
}

func (o ListMuteRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMuteRuleRequest struct{}"
	}

	return strings.Join([]string{"ListMuteRuleRequest", string(data)}, " ")
}
