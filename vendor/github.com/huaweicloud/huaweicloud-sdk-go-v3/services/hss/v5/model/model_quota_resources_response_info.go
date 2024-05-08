package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// QuotaResourcesResponseInfo 配额资源
type QuotaResourcesResponseInfo struct {

	// 主机安全配额的资源ID
	ResourceId *string `json:"resource_id,omitempty"`

	// 资源规格编码，包含如下:   - hss.version.basic : 基础版   - hss.version.advanced : 专业版   - hss.version.enterprise : 企业版   - hss.version.premium : 旗舰版   - hss.version.wtp : 网页防篡改版   - hss.version.container : 容器版
	Version *string `json:"version,omitempty"`

	// 配额状态   - normal : 正常   - expired : 已过期   - freeze : 已冻结
	QuotaStatus *string `json:"quota_status,omitempty"`

	// 使用状态   - idle : 空闲   - used : 使用中
	UsedStatus *string `json:"used_status,omitempty"`

	// 主机ID
	HostId *string `json:"host_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 计费模式   - packet_cycle : 包周期   - on_demand : 按需
	ChargingMode *string `json:"charging_mode,omitempty"`

	// 标签
	Tags *[]TagInfo `json:"tags,omitempty"`

	// 过期时间，-1表示没有到期时间
	ExpireTime *int64 `json:"expire_time,omitempty"`

	// 是否共享配额   - shared：共享的   - unshared：非共享的
	SharedQuota *string `json:"shared_quota,omitempty"`

	// 企业项目ID
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 所属企业项目名称
	EnterpriseProjectName *string `json:"enterprise_project_name,omitempty"`
}

func (o QuotaResourcesResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QuotaResourcesResponseInfo struct{}"
	}

	return strings.Join([]string{"QuotaResourcesResponseInfo", string(data)}, " ")
}
