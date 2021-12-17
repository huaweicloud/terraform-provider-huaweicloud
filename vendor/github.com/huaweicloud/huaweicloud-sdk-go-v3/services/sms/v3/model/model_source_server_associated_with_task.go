package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 任务关联的源端信息
type SourceServerAssociatedWithTask struct {
	// 源端在SMS数据库中的ID

	Id *string `json:"id,omitempty"`
	// 源端服务器ip，注册源端时必选，更新非必选

	Ip string `json:"ip"`
	// 用来区分不同源端服务器的名称

	Name string `json:"name"`
	// 源端服务器的OS类型，分为Windows和Linux，注册必选，更新非必选

	OsType SourceServerAssociatedWithTaskOsType `json:"os_type"`
	// 操作系统版本，注册必选，更新非必选

	OsVersion *string `json:"os_version,omitempty"`
	// 是否是OEM操作系统(Windows)

	OemSystem *bool `json:"oem_system,omitempty"`
	// 当前源端服务器状态

	State *SourceServerAssociatedWithTaskState `json:"state,omitempty"`
}

func (o SourceServerAssociatedWithTask) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SourceServerAssociatedWithTask struct{}"
	}

	return strings.Join([]string{"SourceServerAssociatedWithTask", string(data)}, " ")
}

type SourceServerAssociatedWithTaskOsType struct {
	value string
}

type SourceServerAssociatedWithTaskOsTypeEnum struct {
	WINDOWS SourceServerAssociatedWithTaskOsType
	LINUX   SourceServerAssociatedWithTaskOsType
}

func GetSourceServerAssociatedWithTaskOsTypeEnum() SourceServerAssociatedWithTaskOsTypeEnum {
	return SourceServerAssociatedWithTaskOsTypeEnum{
		WINDOWS: SourceServerAssociatedWithTaskOsType{
			value: "WINDOWS",
		},
		LINUX: SourceServerAssociatedWithTaskOsType{
			value: "LINUX",
		},
	}
}

func (c SourceServerAssociatedWithTaskOsType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourceServerAssociatedWithTaskOsType) UnmarshalJSON(b []byte) error {
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

type SourceServerAssociatedWithTaskState struct {
	value string
}

type SourceServerAssociatedWithTaskStateEnum struct {
	UNAVAILABLE SourceServerAssociatedWithTaskState
	WAITING     SourceServerAssociatedWithTaskState
	INITIALIZE  SourceServerAssociatedWithTaskState
	REPLICATE   SourceServerAssociatedWithTaskState
	SYNCING     SourceServerAssociatedWithTaskState
	STOPPING    SourceServerAssociatedWithTaskState
	STOPPED     SourceServerAssociatedWithTaskState
	DELETING    SourceServerAssociatedWithTaskState
	ERROR       SourceServerAssociatedWithTaskState
	CLONING     SourceServerAssociatedWithTaskState
	TESTING     SourceServerAssociatedWithTaskState
	FINISHED    SourceServerAssociatedWithTaskState
}

func GetSourceServerAssociatedWithTaskStateEnum() SourceServerAssociatedWithTaskStateEnum {
	return SourceServerAssociatedWithTaskStateEnum{
		UNAVAILABLE: SourceServerAssociatedWithTaskState{
			value: "unavailable",
		},
		WAITING: SourceServerAssociatedWithTaskState{
			value: "waiting",
		},
		INITIALIZE: SourceServerAssociatedWithTaskState{
			value: "initialize",
		},
		REPLICATE: SourceServerAssociatedWithTaskState{
			value: "replicate",
		},
		SYNCING: SourceServerAssociatedWithTaskState{
			value: "syncing",
		},
		STOPPING: SourceServerAssociatedWithTaskState{
			value: "stopping",
		},
		STOPPED: SourceServerAssociatedWithTaskState{
			value: "stopped",
		},
		DELETING: SourceServerAssociatedWithTaskState{
			value: "deleting",
		},
		ERROR: SourceServerAssociatedWithTaskState{
			value: "error",
		},
		CLONING: SourceServerAssociatedWithTaskState{
			value: "cloning",
		},
		TESTING: SourceServerAssociatedWithTaskState{
			value: "testing",
		},
		FINISHED: SourceServerAssociatedWithTaskState{
			value: "finished",
		},
	}
}

func (c SourceServerAssociatedWithTaskState) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourceServerAssociatedWithTaskState) UnmarshalJSON(b []byte) error {
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
