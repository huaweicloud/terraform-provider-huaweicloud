package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteGroupResponse Response Object
type BatchDeleteGroupResponse struct {

	// 删除失败的消费组列表。
	FailedGroups *[]BatchDeleteGroupRespFailedGroups `json:"failed_groups,omitempty"`

	// 删除失败的个数
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o BatchDeleteGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteGroupResponse struct{}"
	}

	return strings.Join([]string{"BatchDeleteGroupResponse", string(data)}, " ")
}
