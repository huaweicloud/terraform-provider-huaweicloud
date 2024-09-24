package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShareBackups struct {

	// 备份ID。
	Id *string `json:"id,omitempty"`

	// 备份名字。
	Name *string `json:"name,omitempty"`

	// 备份开始时间。
	BeginTime *string `json:"begin_time,omitempty"`

	// 备份结束时间。
	EndTime *string `json:"end_time,omitempty"`

	// 备份状态，取值：BUILDING：备份中，COMPLETED：备份完成，FAILED：备份失败，DELETING：备份删除中。
	Status *string `json:"status,omitempty"`

	// 备份大小，单位：KB。
	Size *float64 `json:"size,omitempty"`

	// 备份类型，取值：\"auto\"自动全量备份，“manual”手动全量备份。
	Type *string `json:"type,omitempty"`

	// 备份方法。
	BackupMethod *string `json:"backup_method,omitempty"`

	// 备份所在实例ID。
	InstanceId *string `json:"instance_id,omitempty"`

	// 备份所在实例名称。
	InstanceName *string `json:"instance_name,omitempty"`

	// 备份所在实例状态。
	InstanceStatus *string `json:"instance_status,omitempty"`

	// 数据库版本信息。
	Datastore *interface{} `json:"datastore,omitempty"`

	// 共享者用户名称。
	UserName *string `json:"user_name,omitempty"`
}

func (o ShareBackups) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShareBackups struct{}"
	}

	return strings.Join([]string{"ShareBackups", string(data)}, " ")
}
