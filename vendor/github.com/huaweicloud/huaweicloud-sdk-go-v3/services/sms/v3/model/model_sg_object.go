package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 安全组object
type SgObject struct {
	// 安全组ID

	Id string `json:"id"`
	// 安全组名称

	Name string `json:"name"`
}

func (o SgObject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SgObject struct{}"
	}

	return strings.Join([]string{"SgObject", string(data)}, " ")
}
