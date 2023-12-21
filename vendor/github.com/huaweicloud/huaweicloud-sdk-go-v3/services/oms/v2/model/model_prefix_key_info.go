package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PrefixKeyInfo 前缀key
type PrefixKeyInfo struct {

	// 键
	Keys []string `json:"keys"`
}

func (o PrefixKeyInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PrefixKeyInfo struct{}"
	}

	return strings.Join([]string{"PrefixKeyInfo", string(data)}, " ")
}
