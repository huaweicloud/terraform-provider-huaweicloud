package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListErrorServersResponse struct {
	// 迁移过程中发生错误的源端数量

	Count *int32 `json:"count,omitempty"`
	// 迁移过程中发生的错误详情

	MigrationErrors *[]MigrationErrors `json:"migration_errors,omitempty"`
	HttpStatusCode  int                `json:"-"`
}

func (o ListErrorServersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListErrorServersResponse struct{}"
	}

	return strings.Join([]string{"ListErrorServersResponse", string(data)}, " ")
}
