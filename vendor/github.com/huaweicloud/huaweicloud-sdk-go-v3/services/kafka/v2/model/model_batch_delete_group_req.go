package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchDeleteGroupReq struct {

	// 所有需要删除的消费组ID。
	GroupIds []string `json:"group_ids"`
}

func (o BatchDeleteGroupReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteGroupReq struct{}"
	}

	return strings.Join([]string{"BatchDeleteGroupReq", string(data)}, " ")
}
