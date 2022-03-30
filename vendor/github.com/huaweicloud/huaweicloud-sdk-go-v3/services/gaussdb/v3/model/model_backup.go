package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type Backup struct {
	// 备份ID

	Id *string `json:"id,omitempty"`
	// 备份名称。

	Name *string `json:"name,omitempty"`
	// 备份描述。

	Descriprion *string `json:"descriprion,omitempty"`
	// 备份开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”，其中T指时间字段的开始；Z指时区偏移量。

	BeginTime *string `json:"begin_time,omitempty"`
	// 备份状态

	Status *BackupStatus `json:"status,omitempty"`
	// 备份类型，取值：

	Type *BackupType `json:"type,omitempty"`
	// 实例ID。

	InstanceId *string `json:"instance_id,omitempty"`
}

func (o Backup) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Backup struct{}"
	}

	return strings.Join([]string{"Backup", string(data)}, " ")
}

type BackupStatus struct {
	value string
}

type BackupStatusEnum struct {
	BUILDING  BackupStatus
	COMPLETED BackupStatus
	FAILED    BackupStatus
	AVAILABLE BackupStatus
}

func GetBackupStatusEnum() BackupStatusEnum {
	return BackupStatusEnum{
		BUILDING: BackupStatus{
			value: "BUILDING：备份中。",
		},
		COMPLETED: BackupStatus{
			value: "COMPLETED：备份完成。",
		},
		FAILED: BackupStatus{
			value: "FAILED：备份失败。",
		},
		AVAILABLE: BackupStatus{
			value: "AVAILABLE：备份可用。",
		},
	}
}

func (c BackupStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BackupStatus) UnmarshalJSON(b []byte) error {
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

type BackupType struct {
	value string
}

type BackupTypeEnum struct {
	MANUAL BackupType
}

func GetBackupTypeEnum() BackupTypeEnum {
	return BackupTypeEnum{
		MANUAL: BackupType{
			value: "manual：手动全量备份。",
		},
	}
}

func (c BackupType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BackupType) UnmarshalJSON(b []byte) error {
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
