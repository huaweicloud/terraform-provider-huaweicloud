package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListMigprojectsResponse struct {
	// 查询到的迁移项目的数量

	Count *int32 `json:"count,omitempty"`
	// 查询到的迁移项目详情

	Migprojects    *[]MigprojectsResponseBody `json:"migprojects,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o ListMigprojectsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMigprojectsResponse struct{}"
	}

	return strings.Join([]string{"ListMigprojectsResponse", string(data)}, " ")
}
