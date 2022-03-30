package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlResizeFlavor struct {
	// 规格码

	SpecCode string `json:"spec_code"`
}

func (o MysqlResizeFlavor) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlResizeFlavor struct{}"
	}

	return strings.Join([]string{"MysqlResizeFlavor", string(data)}, " ")
}
