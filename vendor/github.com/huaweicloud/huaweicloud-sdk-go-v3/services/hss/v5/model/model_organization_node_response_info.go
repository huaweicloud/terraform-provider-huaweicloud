package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// OrganizationNodeResponseInfo 组织结构树
type OrganizationNodeResponseInfo struct {

	// 父节点Id
	ParentId *string `json:"parent_id,omitempty"`

	// 节点account_id
	Id *string `json:"id,omitempty"`

	// 组织的统一资源名称,格式：organizations::{management_account_id}:xxxxx:{org_id}/xxxxxxxx。
	Urn *string `json:"urn,omitempty"`

	// 名称
	Name *string `json:"name,omitempty"`

	// 节点类型，unit:组织单元、account:账号
	OrgType *string `json:"org_type,omitempty"`

	// 组织或账号是否已授权。   - true: 已授权（无需授权）。   - false: 未授权。
	Delegated *bool `json:"delegated,omitempty"`
}

func (o OrganizationNodeResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OrganizationNodeResponseInfo struct{}"
	}

	return strings.Join([]string{"OrganizationNodeResponseInfo", string(data)}, " ")
}
