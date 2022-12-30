package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 服务器组ID
type GroupId struct {
}

func (o GroupId) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GroupId struct{}"
	}

	return strings.Join([]string{"GroupId", string(data)}, " ")
}
