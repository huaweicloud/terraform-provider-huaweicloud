package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListServersRequest struct {
	// 源端服务器状态

	State *ListServersRequestState `json:"state,omitempty"`
	// 源端服务器名称

	Name *string `json:"name,omitempty"`
	// 源端服务器ID

	Id *string `json:"id,omitempty"`
	// 源端服务器IP地址

	Ip *string `json:"ip,omitempty"`
	// 迁移项目id，填写该参数将查询迁移项目下的所有虚拟机

	Migproject *string `json:"migproject,omitempty"`
	// 每一页记录的源端服务器数量，0表示用默认值 200

	Limit *int32 `json:"limit,omitempty"`
	// 偏移量，默认值0

	Offset *int32 `json:"offset,omitempty"`
	// 根据迁移周期查询

	MigrationCycle *ListServersRequestMigrationCycle `json:"migration_cycle,omitempty"`
	// 查询失去连接的源端

	Connected *bool `json:"connected,omitempty"`
	// 需要查询的企业项目id

	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ListServersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListServersRequest struct{}"
	}

	return strings.Join([]string{"ListServersRequest", string(data)}, " ")
}

type ListServersRequestState struct {
	value string
}

type ListServersRequestStateEnum struct {
	UNAVAILABLE ListServersRequestState
	WAITING     ListServersRequestState
	INITIALIZE  ListServersRequestState
	REPLICATE   ListServersRequestState
	SYNCING     ListServersRequestState
	STOPPING    ListServersRequestState
	STOPPED     ListServersRequestState
	DELETING    ListServersRequestState
	ERROR       ListServersRequestState
	CLONING     ListServersRequestState
	CUTOVERING  ListServersRequestState
	FINISHED    ListServersRequestState
}

func GetListServersRequestStateEnum() ListServersRequestStateEnum {
	return ListServersRequestStateEnum{
		UNAVAILABLE: ListServersRequestState{
			value: "unavailable",
		},
		WAITING: ListServersRequestState{
			value: "waiting",
		},
		INITIALIZE: ListServersRequestState{
			value: "initialize",
		},
		REPLICATE: ListServersRequestState{
			value: "replicate",
		},
		SYNCING: ListServersRequestState{
			value: "syncing",
		},
		STOPPING: ListServersRequestState{
			value: "stopping",
		},
		STOPPED: ListServersRequestState{
			value: "stopped",
		},
		DELETING: ListServersRequestState{
			value: "deleting",
		},
		ERROR: ListServersRequestState{
			value: "error",
		},
		CLONING: ListServersRequestState{
			value: "cloning",
		},
		CUTOVERING: ListServersRequestState{
			value: "cutovering",
		},
		FINISHED: ListServersRequestState{
			value: "finished",
		},
	}
}

func (c ListServersRequestState) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListServersRequestState) UnmarshalJSON(b []byte) error {
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

type ListServersRequestMigrationCycle struct {
	value string
}

type ListServersRequestMigrationCycleEnum struct {
	CHECKING    ListServersRequestMigrationCycle
	SETTING     ListServersRequestMigrationCycle
	REPLICATING ListServersRequestMigrationCycle
	SYNCING     ListServersRequestMigrationCycle
	CUTOVERING  ListServersRequestMigrationCycle
	CUTOVERED   ListServersRequestMigrationCycle
}

func GetListServersRequestMigrationCycleEnum() ListServersRequestMigrationCycleEnum {
	return ListServersRequestMigrationCycleEnum{
		CHECKING: ListServersRequestMigrationCycle{
			value: "checking",
		},
		SETTING: ListServersRequestMigrationCycle{
			value: "setting",
		},
		REPLICATING: ListServersRequestMigrationCycle{
			value: "replicating",
		},
		SYNCING: ListServersRequestMigrationCycle{
			value: "syncing",
		},
		CUTOVERING: ListServersRequestMigrationCycle{
			value: "cutovering",
		},
		CUTOVERED: ListServersRequestMigrationCycle{
			value: "cutovered",
		},
	}
}

func (c ListServersRequestMigrationCycle) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListServersRequestMigrationCycle) UnmarshalJSON(b []byte) error {
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
