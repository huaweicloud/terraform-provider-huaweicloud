package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteStreamForbiddenRequest Request Object
type DeleteStreamForbiddenRequest struct {

	// 推流域名
	Domain string `json:"domain"`

	// RTMP应用名称
	AppName string `json:"app_name"`

	// 流名称
	StreamName string `json:"stream_name"`
}

func (o DeleteStreamForbiddenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteStreamForbiddenRequest struct{}"
	}

	return strings.Join([]string{"DeleteStreamForbiddenRequest", string(data)}, " ")
}
