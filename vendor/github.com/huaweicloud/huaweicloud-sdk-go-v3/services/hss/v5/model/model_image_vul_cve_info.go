package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ImageVulCveInfo cve info
type ImageVulCveInfo struct {

	// cve id
	CveId *string `json:"cve_id,omitempty"`

	// CVSS分数
	CvssScore *float32 `json:"cvss_score,omitempty"`

	// 公布时间，时间单位 毫秒（ms）
	PublishTime *int64 `json:"publish_time,omitempty"`

	// cve描述
	Description *string `json:"description,omitempty"`
}

func (o ImageVulCveInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ImageVulCveInfo struct{}"
	}

	return strings.Join([]string{"ImageVulCveInfo", string(data)}, " ")
}
