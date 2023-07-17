package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateGroupReq struct {

	// 消费组名称
	GroupName string `json:"group_name"`
}

func (o CreateGroupReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateGroupReq struct{}"
	}

	return strings.Join([]string{"CreateGroupReq", string(data)}, " ")
}
