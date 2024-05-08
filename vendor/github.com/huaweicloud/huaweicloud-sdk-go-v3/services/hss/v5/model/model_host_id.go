package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HostId 主机ID
type HostId struct {
}

func (o HostId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HostId struct{}"
	}

	return strings.Join([]string{"HostId", string(data)}, " ")
}
