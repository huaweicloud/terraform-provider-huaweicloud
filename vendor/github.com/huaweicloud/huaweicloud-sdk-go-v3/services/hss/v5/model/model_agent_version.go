package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AgentVersion agent版本
type AgentVersion struct {
}

func (o AgentVersion) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AgentVersion struct{}"
	}

	return strings.Join([]string{"AgentVersion", string(data)}, " ")
}
