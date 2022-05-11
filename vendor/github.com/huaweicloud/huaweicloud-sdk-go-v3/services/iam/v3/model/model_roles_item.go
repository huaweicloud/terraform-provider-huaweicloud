package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RolesItem struct {

	// 权限所在目录。
	Catalog string `json:"catalog"`

	// 权限展示名称。
	DisplayName string `json:"display_name"`

	// 权限的英文描述。
	Description string `json:"description"`

	// 权限的中文描述信息。
	DescriptionCn string `json:"description_cn"`

	// 权限所属账号ID。
	DomainId string `json:"domain_id"`

	// 该参数值为fine_grained时，标识此权限为系统内置的策略。
	Flag string `json:"flag"`

	// 权限Id。
	Id string `json:"id"`

	// 权限名称。
	Name string `json:"name"`

	Policy *RolePolicy `json:"policy"`

	// 权限的显示模式。 > - AX表示在domain层显示。 > - XA表示在project层显示。 > - AA表示在domain和project层均显示。 > - XX表示在domain和project层均不显示。 > - 自定义策略的显示模式只能为AX或者XA，不能在domain层和project层都显示（AA），或者在domain层和project层都不显示（XX）。
	Type string `json:"type"`
}

func (o RolesItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RolesItem struct{}"
	}

	return strings.Join([]string{"RolesItem", string(data)}, " ")
}
