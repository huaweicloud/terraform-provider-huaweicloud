package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 创建任务的参数
type PostTask struct {
	// 任务名称

	Name string `json:"name"`
	// 任务类型

	Type PostTaskType `json:"type"`
	// 迁移后是否启动目的端虚拟机

	StartTargetServer *bool `json:"start_target_server,omitempty"`
	// 操作系统类型

	OsType string `json:"os_type"`

	SourceServer *SourceServerByTask `json:"source_server"`

	TargetServer *TargetServerByTask `json:"target_server"`
	// 迁移ip，如果是自动创建虚拟机，不需要此参数

	MigrationIp *string `json:"migration_ip,omitempty"`
	// region的名称

	RegionName string `json:"region_name"`
	// region id

	RegionId string `json:"region_id"`
	// 项目名称

	ProjectName string `json:"project_name"`
	// 项目id

	ProjectId string `json:"project_id"`
	// 自动创建虚拟机使用模板

	VmTemplateId *string `json:"vm_template_id,omitempty"`
	// 是否使用公网ip

	UsePublicIp *bool `json:"use_public_ip,omitempty"`
	// 复制或者同步后是否会继续持续同步，不添加则默认是false

	Syncing *bool `json:"syncing,omitempty"`
}

func (o PostTask) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostTask struct{}"
	}

	return strings.Join([]string{"PostTask", string(data)}, " ")
}

type PostTaskType struct {
	value string
}

type PostTaskTypeEnum struct {
	MIGRATE_FILE  PostTaskType
	MIGRATE_BLOCK PostTaskType
}

func GetPostTaskTypeEnum() PostTaskTypeEnum {
	return PostTaskTypeEnum{
		MIGRATE_FILE: PostTaskType{
			value: "MIGRATE_FILE",
		},
		MIGRATE_BLOCK: PostTaskType{
			value: "MIGRATE_BLOCK",
		},
	}
}

func (c PostTaskType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PostTaskType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
