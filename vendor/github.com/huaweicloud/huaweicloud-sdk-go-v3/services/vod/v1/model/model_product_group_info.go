package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProductGroupInfo struct {

	// 模板组ID
	GroupId *string `json:"group_id,omitempty"`

	// 产物信息
	Products *[]ProductUrlInfo `json:"products,omitempty"`
}

func (o ProductGroupInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProductGroupInfo struct{}"
	}

	return strings.Join([]string{"ProductGroupInfo", string(data)}, " ")
}
