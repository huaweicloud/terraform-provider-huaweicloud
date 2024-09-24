package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProductDelFailInfo struct {

	// 删除产物的URL
	Url *string `json:"url,omitempty"`

	// 删除产物失败的原因
	FailReason *string `json:"fail_reason,omitempty"`
}

func (o ProductDelFailInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProductDelFailInfo struct{}"
	}

	return strings.Join([]string{"ProductDelFailInfo", string(data)}, " ")
}
