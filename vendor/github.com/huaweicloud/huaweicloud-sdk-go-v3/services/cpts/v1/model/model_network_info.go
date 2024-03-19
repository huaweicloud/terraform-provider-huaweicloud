package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NetworkInfo struct {

	// 网络类型（固定值：internet）
	NetworkType string `json:"network_type"`
}

func (o NetworkInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NetworkInfo struct{}"
	}

	return strings.Join([]string{"NetworkInfo", string(data)}, " ")
}
