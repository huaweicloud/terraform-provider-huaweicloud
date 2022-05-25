package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RecordControlInfo struct {

	// 直播推流域名
	PublishDomain string `json:"publish_domain"`

	// 应用名
	App string `json:"app"`

	// 待启动或停止录制的流名
	Stream string `json:"stream"`
}

func (o RecordControlInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecordControlInfo struct{}"
	}

	return strings.Join([]string{"RecordControlInfo", string(data)}, " ")
}
