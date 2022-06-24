package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowApiInfoRequest struct {

	// 版本信息。
	Version string `json:"version"`
}

func (o ShowApiInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowApiInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowApiInfoRequest", string(data)}, " ")
}
