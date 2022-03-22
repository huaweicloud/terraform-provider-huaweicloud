package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SetQuota struct {
	// 企业项目ID。

	EnterpriseProjectId string `json:"enterprise_project_id"`
	// 实例个数配额。取值范围0~1000。(如果已经存在实例，应该大于已经存在的实例个数)

	InstanceQuota int32 `json:"instance_quota"`
	// CPU核数配额。取值范围0~3600000。(如果已经存在实例，应该大于已经占用的cpu个数)

	VcpusQuota int32 `json:"vcpus_quota"`
	// 内存使用配额，单位为GB。取值范围0~19200000。(如果已经存在实例，应该大于已经占用的内存数)

	RamQuota int32 `json:"ram_quota"`
}

func (o SetQuota) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetQuota struct{}"
	}

	return strings.Join([]string{"SetQuota", string(data)}, " ")
}
