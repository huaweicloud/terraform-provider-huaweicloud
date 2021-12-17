package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 源端服务器列表信息
type SourceServersResponseBody struct {
	// 源端服务器id

	Id *string `json:"id,omitempty"`
	// 源端服务器的ip地址

	Ip *string `json:"ip,omitempty"`
	// 源端服务器名称

	Name *string `json:"name,omitempty"`
	// 企业项目id

	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
	// 源端服务器的注册时间

	AddDate *int64 `json:"add_date,omitempty"`
	// 操作系统类型，OS_TYPE (WINDOWS/LINUX)

	OsType *SourceServersResponseBodyOsType `json:"os_type,omitempty"`
	// 系统详细版本号，如CENTOS7.6等

	OsVersion *string `json:"os_version,omitempty"`
	// 是否是OEM操作系统(Windows)

	OemSystem *bool `json:"oem_system,omitempty"`
	// 源端服务器状态

	State *SourceServersResponseBodyState `json:"state,omitempty"`
	// 源端服务器与主机迁移服务端是否连接

	Connected *bool `json:"connected,omitempty"`
	// 源端CPU核心数

	CpuQuantity *int32 `json:"cpu_quantity,omitempty"`
	// 源端物理内存大小（单位：字节）

	Memory *int64 `json:"memory,omitempty"`

	CurrentTask *TaskByServerSources `json:"current_task,omitempty"`
	// 源端校验检查项列表

	Checks *[]EnvironmentCheck `json:"checks,omitempty"`

	InitTargetServer *InitTargetServer `json:"init_target_server,omitempty"`
	// 已复制的大小（单位：字节）

	Replicatesize *int64 `json:"replicatesize,omitempty"`
	// 迁移周期（migration_cycle）上一次变化的时间

	StageActionTime *int64 `json:"stage_action_time,omitempty"`
	// 需要迁移的数据量总大小（单位：字节）

	Totalsize *int64 `json:"totalsize,omitempty"`
	// Agent上一次连接状态发生变化的时间

	LastVisitTime *int64 `json:"last_visit_time,omitempty"`
	// 迁移周期

	MigrationCycle *SourceServersResponseBodyMigrationCycle `json:"migration_cycle,omitempty"`
	// 源端状态（state）上次发生变化的时间

	StateActionTime *int64 `json:"state_action_time,omitempty"`
}

func (o SourceServersResponseBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SourceServersResponseBody struct{}"
	}

	return strings.Join([]string{"SourceServersResponseBody", string(data)}, " ")
}

type SourceServersResponseBodyOsType struct {
	value string
}

type SourceServersResponseBodyOsTypeEnum struct {
	WINDOWS SourceServersResponseBodyOsType
	LINUX   SourceServersResponseBodyOsType
}

func GetSourceServersResponseBodyOsTypeEnum() SourceServersResponseBodyOsTypeEnum {
	return SourceServersResponseBodyOsTypeEnum{
		WINDOWS: SourceServersResponseBodyOsType{
			value: "WINDOWS",
		},
		LINUX: SourceServersResponseBodyOsType{
			value: "LINUX",
		},
	}
}

func (c SourceServersResponseBodyOsType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourceServersResponseBodyOsType) UnmarshalJSON(b []byte) error {
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

type SourceServersResponseBodyState struct {
	value string
}

type SourceServersResponseBodyStateEnum struct {
	UNAVAILABLE SourceServersResponseBodyState
	WAITING     SourceServersResponseBodyState
	INITIALIZE  SourceServersResponseBodyState
	REPLICATE   SourceServersResponseBodyState
	SYNCING     SourceServersResponseBodyState
	STOPPING    SourceServersResponseBodyState
	STOPPED     SourceServersResponseBodyState
	DELETING    SourceServersResponseBodyState
	ERROR       SourceServersResponseBodyState
	CLONING     SourceServersResponseBodyState
	CUTOVERING  SourceServersResponseBodyState
	FINISHED    SourceServersResponseBodyState
}

func GetSourceServersResponseBodyStateEnum() SourceServersResponseBodyStateEnum {
	return SourceServersResponseBodyStateEnum{
		UNAVAILABLE: SourceServersResponseBodyState{
			value: "unavailable",
		},
		WAITING: SourceServersResponseBodyState{
			value: "waiting",
		},
		INITIALIZE: SourceServersResponseBodyState{
			value: "initialize",
		},
		REPLICATE: SourceServersResponseBodyState{
			value: "replicate",
		},
		SYNCING: SourceServersResponseBodyState{
			value: "syncing",
		},
		STOPPING: SourceServersResponseBodyState{
			value: "stopping",
		},
		STOPPED: SourceServersResponseBodyState{
			value: "stopped",
		},
		DELETING: SourceServersResponseBodyState{
			value: "deleting",
		},
		ERROR: SourceServersResponseBodyState{
			value: "error",
		},
		CLONING: SourceServersResponseBodyState{
			value: "cloning",
		},
		CUTOVERING: SourceServersResponseBodyState{
			value: "cutovering",
		},
		FINISHED: SourceServersResponseBodyState{
			value: "finished",
		},
	}
}

func (c SourceServersResponseBodyState) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourceServersResponseBodyState) UnmarshalJSON(b []byte) error {
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

type SourceServersResponseBodyMigrationCycle struct {
	value string
}

type SourceServersResponseBodyMigrationCycleEnum struct {
	CUTOVERING  SourceServersResponseBodyMigrationCycle
	CUTOVERED   SourceServersResponseBodyMigrationCycle
	CHECKING    SourceServersResponseBodyMigrationCycle
	SETTING     SourceServersResponseBodyMigrationCycle
	REPLICATING SourceServersResponseBodyMigrationCycle
	SYNCING     SourceServersResponseBodyMigrationCycle
}

func GetSourceServersResponseBodyMigrationCycleEnum() SourceServersResponseBodyMigrationCycleEnum {
	return SourceServersResponseBodyMigrationCycleEnum{
		CUTOVERING: SourceServersResponseBodyMigrationCycle{
			value: "cutovering",
		},
		CUTOVERED: SourceServersResponseBodyMigrationCycle{
			value: "cutovered",
		},
		CHECKING: SourceServersResponseBodyMigrationCycle{
			value: "checking",
		},
		SETTING: SourceServersResponseBodyMigrationCycle{
			value: "setting",
		},
		REPLICATING: SourceServersResponseBodyMigrationCycle{
			value: "replicating",
		},
		SYNCING: SourceServersResponseBodyMigrationCycle{
			value: "syncing",
		},
	}
}

func (c SourceServersResponseBodyMigrationCycle) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourceServersResponseBodyMigrationCycle) UnmarshalJSON(b []byte) error {
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
