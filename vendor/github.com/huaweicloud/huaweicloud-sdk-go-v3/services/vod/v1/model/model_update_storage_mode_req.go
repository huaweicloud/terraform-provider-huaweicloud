package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type UpdateStorageModeReq struct {

	// 原媒资id
	AssetId string `json:"asset_id"`

	// 存储模式。 取值如下： - STANDARD：标准存储。 - WARM：低频存储。 - COLD：归档存储。
	StorageMode UpdateStorageModeReqStorageMode `json:"storage_mode"`

	// 归档恢复方式。 取值如下： - TEMP：临时 - FOREVER：永久
	RestoreMode *UpdateStorageModeReqRestoreMode `json:"restore_mode,omitempty"`

	// 从归档转标准存储且选择临时恢复时临时恢复时间。
	Days *int32 `json:"days,omitempty"`

	// 归档恢复选项，快速恢复EXPEDITED，标准恢复STANDARD，默认快速恢复
	RestoreTier *UpdateStorageModeReqRestoreTier `json:"restore_tier,omitempty"`
}

func (o UpdateStorageModeReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateStorageModeReq struct{}"
	}

	return strings.Join([]string{"UpdateStorageModeReq", string(data)}, " ")
}

type UpdateStorageModeReqStorageMode struct {
	value string
}

type UpdateStorageModeReqStorageModeEnum struct {
	STANDARD UpdateStorageModeReqStorageMode
	WARM     UpdateStorageModeReqStorageMode
	COLD     UpdateStorageModeReqStorageMode
}

func GetUpdateStorageModeReqStorageModeEnum() UpdateStorageModeReqStorageModeEnum {
	return UpdateStorageModeReqStorageModeEnum{
		STANDARD: UpdateStorageModeReqStorageMode{
			value: "STANDARD",
		},
		WARM: UpdateStorageModeReqStorageMode{
			value: "WARM",
		},
		COLD: UpdateStorageModeReqStorageMode{
			value: "COLD",
		},
	}
}

func (c UpdateStorageModeReqStorageMode) Value() string {
	return c.value
}

func (c UpdateStorageModeReqStorageMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateStorageModeReqStorageMode) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type UpdateStorageModeReqRestoreMode struct {
	value string
}

type UpdateStorageModeReqRestoreModeEnum struct {
	TEMP    UpdateStorageModeReqRestoreMode
	FOREVER UpdateStorageModeReqRestoreMode
}

func GetUpdateStorageModeReqRestoreModeEnum() UpdateStorageModeReqRestoreModeEnum {
	return UpdateStorageModeReqRestoreModeEnum{
		TEMP: UpdateStorageModeReqRestoreMode{
			value: "TEMP",
		},
		FOREVER: UpdateStorageModeReqRestoreMode{
			value: "FOREVER",
		},
	}
}

func (c UpdateStorageModeReqRestoreMode) Value() string {
	return c.value
}

func (c UpdateStorageModeReqRestoreMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateStorageModeReqRestoreMode) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type UpdateStorageModeReqRestoreTier struct {
	value string
}

type UpdateStorageModeReqRestoreTierEnum struct {
	EXPEDITED UpdateStorageModeReqRestoreTier
	STANDARD  UpdateStorageModeReqRestoreTier
}

func GetUpdateStorageModeReqRestoreTierEnum() UpdateStorageModeReqRestoreTierEnum {
	return UpdateStorageModeReqRestoreTierEnum{
		EXPEDITED: UpdateStorageModeReqRestoreTier{
			value: "EXPEDITED",
		},
		STANDARD: UpdateStorageModeReqRestoreTier{
			value: "STANDARD",
		},
	}
}

func (c UpdateStorageModeReqRestoreTier) Value() string {
	return c.value
}

func (c UpdateStorageModeReqRestoreTier) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateStorageModeReqRestoreTier) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
