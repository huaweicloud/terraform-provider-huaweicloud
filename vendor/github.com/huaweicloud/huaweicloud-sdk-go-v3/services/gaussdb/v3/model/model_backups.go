package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type Backups struct {
	// 备份ID。

	Id *string `json:"id,omitempty"`
	// 备份名称。

	Name *string `json:"name,omitempty"`
	// 备份开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。

	BeginTime *string `json:"begin_time,omitempty"`
	// 备份结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。

	EndTime *string `json:"end_time,omitempty"`
	// 备份状态

	Status *BackupsStatus `json:"status,omitempty"`
	// 备份花费时间(单位：minutes)

	TakeUpTime *int32 `json:"take_up_time,omitempty"`
	// 备份类型

	Type *BackupsType `json:"type,omitempty"`
	// 备份大小，(单位：MB)

	Size *int64 `json:"size,omitempty"`

	Datastore *MysqlDatastore `json:"datastore,omitempty"`
	// 实例ID。

	InstanceId *string `json:"instance_id,omitempty"`
	// 备份级别。当开启一级备份开关时，返回该参数。

	BackupLevel *BackupsBackupLevel `json:"backup_level,omitempty"`
	// 备份文件描述信息

	Description *string `json:"description,omitempty"`
}

func (o Backups) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Backups struct{}"
	}

	return strings.Join([]string{"Backups", string(data)}, " ")
}

type BackupsStatus struct {
	value string
}

type BackupsStatusEnum struct {
	BUILDING  BackupsStatus
	COMPLETED BackupsStatus
	FAILED    BackupsStatus
	AVAILABLE BackupsStatus
}

func GetBackupsStatusEnum() BackupsStatusEnum {
	return BackupsStatusEnum{
		BUILDING: BackupsStatus{
			value: "BUILDING：备份中。",
		},
		COMPLETED: BackupsStatus{
			value: "COMPLETED：备份完成。",
		},
		FAILED: BackupsStatus{
			value: "FAILED：备份失败。",
		},
		AVAILABLE: BackupsStatus{
			value: "AVAILABLE：备份可用。",
		},
	}
}

func (c BackupsStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BackupsStatus) UnmarshalJSON(b []byte) error {
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

type BackupsType struct {
	value string
}

type BackupsTypeEnum struct {
	AUTO   BackupsType
	MANUAL BackupsType
}

func GetBackupsTypeEnum() BackupsTypeEnum {
	return BackupsTypeEnum{
		AUTO: BackupsType{
			value: "auto：自动全量备份。",
		},
		MANUAL: BackupsType{
			value: "manual：手动全量备份。",
		},
	}
}

func (c BackupsType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BackupsType) UnmarshalJSON(b []byte) error {
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

type BackupsBackupLevel struct {
	value string
}

type BackupsBackupLevelEnum struct {
	E_0 BackupsBackupLevel
	E_1 BackupsBackupLevel
	E_2 BackupsBackupLevel
}

func GetBackupsBackupLevelEnum() BackupsBackupLevelEnum {
	return BackupsBackupLevelEnum{
		E_0: BackupsBackupLevel{
			value: "0：备份正在创建中或者备份失败。",
		},
		E_1: BackupsBackupLevel{
			value: "1：一级备份。",
		},
		E_2: BackupsBackupLevel{
			value: "2：二级备份。",
		},
	}
}

func (c BackupsBackupLevel) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BackupsBackupLevel) UnmarshalJSON(b []byte) error {
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
