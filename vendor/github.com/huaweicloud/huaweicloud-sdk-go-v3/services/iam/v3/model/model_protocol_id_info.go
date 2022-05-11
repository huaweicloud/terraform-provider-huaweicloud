package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 协议信息。
type ProtocolIdInfo struct {

	// 协议id。
	Id string `json:"id"`
}

func (o ProtocolIdInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProtocolIdInfo struct{}"
	}

	return strings.Join([]string{"ProtocolIdInfo", string(data)}, " ")
}
