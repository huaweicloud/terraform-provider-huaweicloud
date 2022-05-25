package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteStreamForbiddenRequest struct {

	// op账号需要携带的特定project_id，当使用op账号时该值为所操作租户的project_id
	SpecifyProject *string `json:"specify_project,omitempty"`

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
