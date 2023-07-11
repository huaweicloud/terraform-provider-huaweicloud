package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PolicyGroupId 策略组ID
type PolicyGroupId struct {
}

func (o PolicyGroupId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PolicyGroupId struct{}"
	}

	return strings.Join([]string{"PolicyGroupId", string(data)}, " ")
}
