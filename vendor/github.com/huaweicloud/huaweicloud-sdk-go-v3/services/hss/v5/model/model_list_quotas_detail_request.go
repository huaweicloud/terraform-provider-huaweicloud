package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListQuotasDetailRequest Request Object
type ListQuotasDetailRequest struct {

	// Region ID
	Region *string `json:"region,omitempty"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 主机开通的版本，包含如下7种输入。   - hss.version.null ：无。   - hss.version.basic ：基础版。   - hss.version.advanced ：专业版。   - hss.version.enterprise ：企业版。   - hss.version.premium ：旗舰版。   - hss.version.wtp ：网页防篡改版。   - hss.version.container.enterprise：容器版。
	Version *string `json:"version,omitempty"`

	// 类别，包含如下几种：   - host_resource ：HOST_RESOURCE   - container_resource ：CONTAINER_RESOURCE
	Category *string `json:"category,omitempty"`

	// 配额状态，包含如下几种：   - normal ： QUOTA_STATUS_NORMAL   - expired ：QUOTA_STATUS_EXPIRED   - freeze ：QUOTA_STATUS_FREEZE
	QuotaStatus *string `json:"quota_status,omitempty"`

	// 使用状态，包含如下几种：   - idle ：USED_STATUS_IDLE   - used ：USED_STATUS_USED
	UsedStatus *string `json:"used_status,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// HSS配额的资源ID
	ResourceId *string `json:"resource_id,omitempty"`

	// 收费模式，包含如下2种。   - packet_cycle ：包年/包月。   - on_demand ：按需。
	ChargingMode *string `json:"charging_mode,omitempty"`

	// 每页数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListQuotasDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListQuotasDetailRequest struct{}"
	}

	return strings.Join([]string{"ListQuotasDetailRequest", string(data)}, " ")
}
