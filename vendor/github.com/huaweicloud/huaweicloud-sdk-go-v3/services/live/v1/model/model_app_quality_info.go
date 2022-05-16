package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AppQualityInfo struct {

	// 应用名称
	AppName *string `json:"app_name,omitempty"`

	// 视频质量信息
	QualityInfo *[]QualityInfo `json:"quality_info,omitempty"`
}

func (o AppQualityInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AppQualityInfo struct{}"
	}

	return strings.Join([]string{"AppQualityInfo", string(data)}, " ")
}
