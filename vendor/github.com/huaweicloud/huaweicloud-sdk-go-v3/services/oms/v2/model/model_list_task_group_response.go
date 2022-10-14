package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListTaskGroupResponse struct {

	// 满足查询条件的任务组总数
	Count *int64 `json:"count,omitempty"`

	// 查询的迁移任务组详情
	Taskgroups     *[]TaskGroupResp `json:"taskgroups,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ListTaskGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTaskGroupResponse struct{}"
	}

	return strings.Join([]string{"ListTaskGroupResponse", string(data)}, " ")
}
