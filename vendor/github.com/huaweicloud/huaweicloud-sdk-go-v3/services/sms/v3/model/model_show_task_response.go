package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Response Object
type ShowTaskResponse struct {
	// 任务名称（用户自定义）

	Name *string `json:"name,omitempty"`
	// 任务类型，创建时必选，更新时可选

	Type *ShowTaskResponseType `json:"type,omitempty"`
	// 操作系统类型，分为WINDOWS和LINUX，创建时必选，更新时可选

	OsType *ShowTaskResponseOsType `json:"os_type,omitempty"`
	// 迁移任务id

	Id *string `json:"id,omitempty"`
	// 进程优先级  0：低  1：标准（默认）  2：高

	Priority *ShowTaskResponsePriority `json:"priority,omitempty"`
	// 目的端服务器的区域ID

	RegionId *string `json:"region_id,omitempty"`
	// 迁移完成后是否启动目的端服务器  true：启动  false：停止

	StartTargetServer *bool `json:"start_target_server,omitempty"`
	// 企业项目id

	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
	// 目的端服务器的IP地址。  公网迁移时请填写弹性IP地址  专线迁移时请填写私有IP地址

	MigrationIp *string `json:"migration_ip,omitempty"`
	// 目的端服务器的区域名称

	RegionName *string `json:"region_name,omitempty"`
	// 目的端服务器所在项目名称

	ProjectName *string `json:"project_name,omitempty"`
	// 目的端服务器所在项目ID

	ProjectId *string `json:"project_id,omitempty"`
	// 模板ID

	VmTemplateId *string `json:"vm_template_id,omitempty"`

	SourceServer *SourceServer `json:"source_server,omitempty"`

	TargetServer *TaskTargetServer `json:"target_server,omitempty"`
	// 任务状态

	State *string `json:"state,omitempty"`
	// 预估完成时间

	EstimateCompleteTime *int64 `json:"estimate_complete_time,omitempty"`
	// 连接状态

	Connected *bool `json:"connected,omitempty"`
	// 任务创建时间

	CreateDate *int64 `json:"create_date,omitempty"`
	// 任务开始时间

	StartDate *int64 `json:"start_date,omitempty"`
	// 任务结束时间

	FinishDate *int64 `json:"finish_date,omitempty"`
	// 迁移速率，单位：MB/S

	MigrateSpeed *float64 `json:"migrate_speed,omitempty"`
	// 压缩率

	CompressRate *float64 `json:"compress_rate,omitempty"`
	// 错误信息

	ErrorJson *string `json:"error_json,omitempty"`
	// 任务总耗时

	TotalTime *int64 `json:"total_time,omitempty"`
	// 暂时保留float,兼容现网老版本的SMS-Agent

	FloatIp *string `json:"float_ip,omitempty"`
	// 迁移剩余时间（秒）

	RemainSeconds *int64 `json:"remain_seconds,omitempty"`
	// 目的端的快照id

	TargetSnapshotId *string `json:"target_snapshot_id,omitempty"`

	CloneServer *CloneServer `json:"clone_server,omitempty"`
	// 任务包含的子任务列表

	SubTasks       *[]SubTask `json:"sub_tasks,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o ShowTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTaskResponse struct{}"
	}

	return strings.Join([]string{"ShowTaskResponse", string(data)}, " ")
}

type ShowTaskResponseType struct {
	value string
}

type ShowTaskResponseTypeEnum struct {
	MIGRATE_FILE  ShowTaskResponseType
	MIGRATE_BLOCK ShowTaskResponseType
}

func GetShowTaskResponseTypeEnum() ShowTaskResponseTypeEnum {
	return ShowTaskResponseTypeEnum{
		MIGRATE_FILE: ShowTaskResponseType{
			value: "MIGRATE_FILE",
		},
		MIGRATE_BLOCK: ShowTaskResponseType{
			value: "MIGRATE_BLOCK",
		},
	}
}

func (c ShowTaskResponseType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTaskResponseType) UnmarshalJSON(b []byte) error {
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

type ShowTaskResponseOsType struct {
	value string
}

type ShowTaskResponseOsTypeEnum struct {
	WINDOWS ShowTaskResponseOsType
	LINUX   ShowTaskResponseOsType
}

func GetShowTaskResponseOsTypeEnum() ShowTaskResponseOsTypeEnum {
	return ShowTaskResponseOsTypeEnum{
		WINDOWS: ShowTaskResponseOsType{
			value: "WINDOWS",
		},
		LINUX: ShowTaskResponseOsType{
			value: "LINUX",
		},
	}
}

func (c ShowTaskResponseOsType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTaskResponseOsType) UnmarshalJSON(b []byte) error {
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

type ShowTaskResponsePriority struct {
	value int32
}

type ShowTaskResponsePriorityEnum struct {
	E_0 ShowTaskResponsePriority
	E_1 ShowTaskResponsePriority
	E_2 ShowTaskResponsePriority
}

func GetShowTaskResponsePriorityEnum() ShowTaskResponsePriorityEnum {
	return ShowTaskResponsePriorityEnum{
		E_0: ShowTaskResponsePriority{
			value: 0,
		}, E_1: ShowTaskResponsePriority{
			value: 1,
		}, E_2: ShowTaskResponsePriority{
			value: 2,
		},
	}
}

func (c ShowTaskResponsePriority) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTaskResponsePriority) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("int32")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(int32)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to int32 error")
	}
}
