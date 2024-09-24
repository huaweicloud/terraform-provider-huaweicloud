package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProductGroupDelFailInfo struct {

	// 模板组ID
	GroupId *string `json:"group_id,omitempty"`

	// 模板组删除失败的原因
	FailReason *string `json:"fail_reason,omitempty"`

	// 删除失败的产物的信息
	Products *[]ProductDelFailInfo `json:"products,omitempty"`
}

func (o ProductGroupDelFailInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProductGroupDelFailInfo struct{}"
	}

	return strings.Join([]string{"ProductGroupDelFailInfo", string(data)}, " ")
}
