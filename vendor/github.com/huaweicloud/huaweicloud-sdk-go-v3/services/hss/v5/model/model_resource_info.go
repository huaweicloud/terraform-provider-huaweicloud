package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type ResourceInfo struct {

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 历史开启备份状态，通过筛选可用服务器的error_message或者status判断，如果error_message为空，则没有开启备份，该字段为closed；若不为空，则为opened
	HistoryBackupStatus *ResourceInfoHistoryBackupStatus `json:"history_backup_status,omitempty"`
}

func (o ResourceInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceInfo struct{}"
	}

	return strings.Join([]string{"ResourceInfo", string(data)}, " ")
}

type ResourceInfoHistoryBackupStatus struct {
	value string
}

type ResourceInfoHistoryBackupStatusEnum struct {
	OPENED ResourceInfoHistoryBackupStatus
	CLOSED ResourceInfoHistoryBackupStatus
}

func GetResourceInfoHistoryBackupStatusEnum() ResourceInfoHistoryBackupStatusEnum {
	return ResourceInfoHistoryBackupStatusEnum{
		OPENED: ResourceInfoHistoryBackupStatus{
			value: "opened",
		},
		CLOSED: ResourceInfoHistoryBackupStatus{
			value: "closed",
		},
	}
}

func (c ResourceInfoHistoryBackupStatus) Value() string {
	return c.value
}

func (c ResourceInfoHistoryBackupStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ResourceInfoHistoryBackupStatus) UnmarshalJSON(b []byte) error {
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
