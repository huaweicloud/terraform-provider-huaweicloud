package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchDeleteGroupRespFailedGroups struct {

	// 删除失败的消费组ID。
	GroupId *string `json:"group_id,omitempty"`

	// 删除失败的原因。
	ErrorMessage *string `json:"error_message,omitempty"`
}

func (o BatchDeleteGroupRespFailedGroups) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteGroupRespFailedGroups struct{}"
	}

	return strings.Join([]string{"BatchDeleteGroupRespFailedGroups", string(data)}, " ")
}
