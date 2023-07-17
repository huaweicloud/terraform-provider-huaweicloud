package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AgentId Agent ID
type AgentId struct {
}

func (o AgentId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AgentId struct{}"
	}

	return strings.Join([]string{"AgentId", string(data)}, " ")
}
