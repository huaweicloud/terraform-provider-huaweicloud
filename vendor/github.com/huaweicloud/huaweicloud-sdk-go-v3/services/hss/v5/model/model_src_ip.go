package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SrcIp 攻击源IP
type SrcIp struct {
}

func (o SrcIp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SrcIp struct{}"
	}

	return strings.Join([]string{"SrcIp", string(data)}, " ")
}
