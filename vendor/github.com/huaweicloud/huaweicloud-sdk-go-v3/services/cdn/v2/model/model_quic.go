package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Quic QUIC协议。
type Quic struct {

	// 状态，on：打开，off：关闭。
	Status string `json:"status"`
}

func (o Quic) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Quic struct{}"
	}

	return strings.Join([]string{"Quic", string(data)}, " ")
}
