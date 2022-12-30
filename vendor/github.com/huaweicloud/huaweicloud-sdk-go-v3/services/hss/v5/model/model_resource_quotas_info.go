package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResourceQuotasInfo struct {

	// 主机开通的版本，包含如下6种输入。   - hss.version.null ：无。   - hss.version.basic ：基础版。   - hss.version.enterprise ：企业版。   - hss.version.premium ：旗舰版。   - hss.version.wtp ：网页防篡改版。   - hss.version.container.enterprise：容器版。
	Version *string `json:"version,omitempty"`

	// 总配额数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 已使用配额数
	UsedNum *int32 `json:"used_num,omitempty"`

	// 总配额数
	AvailableNum *int32 `json:"available_num,omitempty"`

	// 可用资源列表
	AvailableResourcesList *[]AvailableResourceIdsInfo `json:"available_resources_list,omitempty"`
}

func (o ResourceQuotasInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceQuotasInfo struct{}"
	}

	return strings.Join([]string{"ResourceQuotasInfo", string(data)}, " ")
}
