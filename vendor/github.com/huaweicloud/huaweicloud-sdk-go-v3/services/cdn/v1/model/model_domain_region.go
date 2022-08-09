package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DomainRegion struct {

	// 域名
	DomainName *string `json:"domain_name,omitempty"`

	// 指标统计数据列表，如果该时间段内无值，则为空数组[]
	RegionIspDetails *[]map[string]interface{} `json:"region_isp_details,omitempty"`
}

func (o DomainRegion) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DomainRegion struct{}"
	}

	return strings.Join([]string{"DomainRegion", string(data)}, " ")
}
