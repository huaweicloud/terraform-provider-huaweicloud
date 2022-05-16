package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StreamForbiddenList struct {

	// 流应用名称
	AppName string `json:"app_name"`

	// 流名称
	StreamName string `json:"stream_name"`

	// 恢复流时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间，不指定则默认7天，最大禁推为90天
	ResumeTime *string `json:"resume_time,omitempty"`
}

func (o StreamForbiddenList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StreamForbiddenList struct{}"
	}

	return strings.Join([]string{"StreamForbiddenList", string(data)}, " ")
}
