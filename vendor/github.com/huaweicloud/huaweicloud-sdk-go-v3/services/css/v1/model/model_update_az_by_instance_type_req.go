package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAzByInstanceTypeReq 切换AZ详情信息。
type UpdateAzByInstanceTypeReq struct {

	// 节点当前所在AZ。
	SourceAz string `json:"source_az"`

	// 期望节点最终分布AZ。
	TargetAz string `json:"target_az"`

	// AZ迁移方式： - multi_az_change：高可用改造。 - az_migrate： AZ平移
	MigrateType string `json:"migrate_type"`

	// 委托名称，委托给CSS，允许CSS调用您的其他云服务。
	Agency string `json:"agency"`

	// 是否进行全量索引快照备份检测。 true：进行全量索引快照备份检测。 false：不进行全量索引快照备份检测。
	IndicesBackupCheck *bool `json:"indices_backup_check,omitempty"`
}

func (o UpdateAzByInstanceTypeReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAzByInstanceTypeReq struct{}"
	}

	return strings.Join([]string{"UpdateAzByInstanceTypeReq", string(data)}, " ")
}
