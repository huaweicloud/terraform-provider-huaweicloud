package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ServerlessInfo struct {

	// Serverless型实例的最小算力，单位RCU，范围0.5~8，步进0.5。
	MinCap string `json:"min_cap"`

	// Serverless型实例的最大算力，单位RCU，范围0.5~8，步进0.5。
	MaxCap string `json:"max_cap"`
}

func (o ServerlessInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ServerlessInfo struct{}"
	}

	return strings.Join([]string{"ServerlessInfo", string(data)}, " ")
}
