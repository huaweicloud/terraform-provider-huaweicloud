package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Response Object
type ShowServerResponse struct {
	// 源端服务器id

	Id *string `json:"id,omitempty"`
	// 源端服务器的ip

	Ip *string `json:"ip,omitempty"`
	// 用来区分不同源端服务器的名称

	Name *string `json:"name,omitempty"`
	// 源端主机名，注册源端必选，更新非必选

	Hostname *string `json:"hostname,omitempty"`
	// 企业项目id

	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
	// 源端服务器注册的时间

	AddDate *int64 `json:"add_date,omitempty"`
	// 源端服务器的OS类型，分为Windows和Linux，注册必选，更新非必选

	OsType *string `json:"os_type,omitempty"`
	// 操作系统版本，注册必选，更新非必选

	OsVersion *string `json:"os_version,omitempty"`
	// 是否是OEM操作系统(Windows)

	OemSystem *bool `json:"oem_system,omitempty"`
	// 当前源端服务器状态

	State *ShowServerResponseState `json:"state,omitempty"`
	// 与Agent连接状态

	Connected *bool `json:"connected,omitempty"`
	// 源端服务器启动类型，如BIOS或者UEFI

	Firmware *ShowServerResponseFirmware `json:"firmware,omitempty"`

	InitTargetServer *InitTargetServer `json:"init_target_server,omitempty"`
	// 源端CPU核心数

	CpuQuantity *int32 `json:"cpu_quantity,omitempty"`
	// 源端服务器物理内存大小，单位MB

	Memory *int64 `json:"memory,omitempty"`

	CurrentTask *TaskByServerSource `json:"current_task,omitempty"`
	// 源端服务器磁盘信息

	Disks *[]ServerDisk `json:"disks,omitempty"`
	// 源端服务器的卷组信息，Linux必选，如果没有卷组，输入[]

	VolumeGroups *[]VolumeGroups `json:"volume_groups,omitempty"`
	// Linux 必选，源端的Btrfs信息。如果源端不存在Btrfs，则为[]

	BtrfsList *[]BtrfsFileSystem `json:"btrfs_list,omitempty"`
	// 源端服务器的网卡信息

	Networks *[]NetWork `json:"networks,omitempty"`
	// 源端环境校验信息

	Checks *[]EnvironmentCheck `json:"checks,omitempty"`
	// 迁移周期

	MigrationCycle *ShowServerResponseMigrationCycle `json:"migration_cycle,omitempty"`
	// 源端状态（state）上次发生变化的时间戳

	StateActionTime *int64 `json:"state_action_time,omitempty"`
	// 已经完成迁移的大小（B）

	Replicatesize *int64 `json:"replicatesize,omitempty"`
	// 需要迁移的数据量总大小（B）

	Totalsize *int64 `json:"totalsize,omitempty"`
	// agent上一次连接状态发生变化的时间戳

	LastVisitTime *int64 `json:"last_visit_time,omitempty"`
	// 迁移周期（migration_cycle）上一次变化的时间戳

	StageActionTime *int64 `json:"stage_action_time,omitempty"`
	// Agent版本信息

	AgentVersion   *string `json:"agent_version,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowServerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowServerResponse struct{}"
	}

	return strings.Join([]string{"ShowServerResponse", string(data)}, " ")
}

type ShowServerResponseState struct {
	value string
}

type ShowServerResponseStateEnum struct {
	UNAVAILABLE ShowServerResponseState
	WAITING     ShowServerResponseState
	INITIALIZE  ShowServerResponseState
	REPLICATE   ShowServerResponseState
	SYNCING     ShowServerResponseState
	STOPPING    ShowServerResponseState
	STOPPED     ShowServerResponseState
	DELETING    ShowServerResponseState
	ERROR       ShowServerResponseState
	CLONING     ShowServerResponseState
	TESTING     ShowServerResponseState
	FINISHED    ShowServerResponseState
}

func GetShowServerResponseStateEnum() ShowServerResponseStateEnum {
	return ShowServerResponseStateEnum{
		UNAVAILABLE: ShowServerResponseState{
			value: "unavailable",
		},
		WAITING: ShowServerResponseState{
			value: "waiting",
		},
		INITIALIZE: ShowServerResponseState{
			value: "initialize",
		},
		REPLICATE: ShowServerResponseState{
			value: "replicate",
		},
		SYNCING: ShowServerResponseState{
			value: "syncing",
		},
		STOPPING: ShowServerResponseState{
			value: "stopping",
		},
		STOPPED: ShowServerResponseState{
			value: "stopped",
		},
		DELETING: ShowServerResponseState{
			value: "deleting",
		},
		ERROR: ShowServerResponseState{
			value: "error",
		},
		CLONING: ShowServerResponseState{
			value: "cloning",
		},
		TESTING: ShowServerResponseState{
			value: "testing",
		},
		FINISHED: ShowServerResponseState{
			value: "finished",
		},
	}
}

func (c ShowServerResponseState) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowServerResponseState) UnmarshalJSON(b []byte) error {
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

type ShowServerResponseFirmware struct {
	value string
}

type ShowServerResponseFirmwareEnum struct {
	BIOS ShowServerResponseFirmware
	UEFI ShowServerResponseFirmware
}

func GetShowServerResponseFirmwareEnum() ShowServerResponseFirmwareEnum {
	return ShowServerResponseFirmwareEnum{
		BIOS: ShowServerResponseFirmware{
			value: "BIOS",
		},
		UEFI: ShowServerResponseFirmware{
			value: "UEFI",
		},
	}
}

func (c ShowServerResponseFirmware) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowServerResponseFirmware) UnmarshalJSON(b []byte) error {
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

type ShowServerResponseMigrationCycle struct {
	value string
}

type ShowServerResponseMigrationCycleEnum struct {
	CUTOVERING  ShowServerResponseMigrationCycle
	CUTOVERED   ShowServerResponseMigrationCycle
	CHECKING    ShowServerResponseMigrationCycle
	SETTING     ShowServerResponseMigrationCycle
	REPLICATING ShowServerResponseMigrationCycle
	SYNCING     ShowServerResponseMigrationCycle
}

func GetShowServerResponseMigrationCycleEnum() ShowServerResponseMigrationCycleEnum {
	return ShowServerResponseMigrationCycleEnum{
		CUTOVERING: ShowServerResponseMigrationCycle{
			value: "cutovering",
		},
		CUTOVERED: ShowServerResponseMigrationCycle{
			value: "cutovered",
		},
		CHECKING: ShowServerResponseMigrationCycle{
			value: "checking",
		},
		SETTING: ShowServerResponseMigrationCycle{
			value: "setting",
		},
		REPLICATING: ShowServerResponseMigrationCycle{
			value: "replicating",
		},
		SYNCING: ShowServerResponseMigrationCycle{
			value: "syncing",
		},
	}
}

func (c ShowServerResponseMigrationCycle) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowServerResponseMigrationCycle) UnmarshalJSON(b []byte) error {
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
