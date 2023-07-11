package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// OccurTime 发生时间，毫秒
type OccurTime struct {
}

func (o OccurTime) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OccurTime struct{}"
	}

	return strings.Join([]string{"OccurTime", string(data)}, " ")
}
