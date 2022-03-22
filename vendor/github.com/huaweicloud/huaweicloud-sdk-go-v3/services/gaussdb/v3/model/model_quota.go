package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Quota struct {
	// 企业项目ID。

	EnterpriseProjectId string `json:"enterprise_project_id"`
	// 企业项目名称。

	EnterpriseProjectName string `json:"enterprise_project_name"`
	// 实例个数配额。

	InstanceQuota int32 `json:"instance_quota"`
	// CPU核数配额。

	VcpusQuota int32 `json:"vcpus_quota"`
	// 内存使用配额，单位为GB。

	RamQuota int32 `json:"ram_quota"`
	// 实例剩余配额。

	AvailabilityInstanceQuota int32 `json:"availability_instance_quota"`
	// CPU核数剩余配额。

	AvailabilityVcpusQuota *int32 `json:"availability_vcpus_quota,omitempty"`
	// 内存剩余配额。

	AvailabilityRamQuota int32 `json:"availability_ram_quota"`
}

func (o Quota) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Quota struct{}"
	}

	return strings.Join([]string{"Quota", string(data)}, " ")
}
