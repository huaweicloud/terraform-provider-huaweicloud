package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PermissionInfo permission列表。
type PermissionInfo struct {

	// permission的ID。
	Id *string `json:"id,omitempty"`

	// permission详情。
	Permission *string `json:"permission,omitempty"`

	// 终端节点服务白名单类型。
	PermissionType *string `json:"permission_type,omitempty"`

	// 白名单的添加时间。
	CreatedAt *string `json:"created_at,omitempty"`
}

func (o PermissionInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PermissionInfo struct{}"
	}

	return strings.Join([]string{"PermissionInfo", string(data)}, " ")
}
