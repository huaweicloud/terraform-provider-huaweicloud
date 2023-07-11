package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateGroupResp 创建结果
type CreateGroupResp struct {
}

func (o CreateGroupResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateGroupResp struct{}"
	}

	return strings.Join([]string{"CreateGroupResp", string(data)}, " ")
}
