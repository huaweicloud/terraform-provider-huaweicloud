package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSqlserverDbUsersRequest Request Object
type ListSqlserverDbUsersRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 分页页码，从1开始。
	Page int32 `json:"page"`

	// 每页数据条数。取值范围[1, 100]。
	Limit int32 `json:"limit"`
}

func (o ListSqlserverDbUsersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSqlserverDbUsersRequest struct{}"
	}

	return strings.Join([]string{"ListSqlserverDbUsersRequest", string(data)}, " ")
}
