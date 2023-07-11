package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TrustProcessInfo struct {

	// 进程路径
	Path *string `json:"path,omitempty"`

	// 进程hash
	Hash *string `json:"hash,omitempty"`
}

func (o TrustProcessInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TrustProcessInfo struct{}"
	}

	return strings.Join([]string{"TrustProcessInfo", string(data)}, " ")
}
