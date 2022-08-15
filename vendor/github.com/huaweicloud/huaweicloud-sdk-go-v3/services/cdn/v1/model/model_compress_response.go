package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CompressResponse struct {

	// GZIP压缩开关。0关闭。1打开
	CompressSwitch int32 `json:"compress_switch"`

	// GZIP压缩规则
	CompressRules *[]CompressRules `json:"compress_rules,omitempty"`
}

func (o CompressResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CompressResponse struct{}"
	}

	return strings.Join([]string{"CompressResponse", string(data)}, " ")
}
