package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ModifyOttChannelState OTT频道修改状态
type ModifyOttChannelState struct {

	// 频道推流域名
	Domain string `json:"domain"`

	// 组名或应用名
	AppName string `json:"app_name"`

	// 频道ID。频道唯一标识，为必填项。频道ID不建议输入下划线“_”，否则会影响转码和截图任务
	Id string `json:"id"`

	// 频道状态 - ON：频道下发成功后，自动启动拉流、转码、录制等功能 - OFF：仅保存频道信息，不启动频道
	State ModifyOttChannelStateState `json:"state"`
}

func (o ModifyOttChannelState) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyOttChannelState struct{}"
	}

	return strings.Join([]string{"ModifyOttChannelState", string(data)}, " ")
}

type ModifyOttChannelStateState struct {
	value string
}

type ModifyOttChannelStateStateEnum struct {
	ON  ModifyOttChannelStateState
	OFF ModifyOttChannelStateState
}

func GetModifyOttChannelStateStateEnum() ModifyOttChannelStateStateEnum {
	return ModifyOttChannelStateStateEnum{
		ON: ModifyOttChannelStateState{
			value: "ON",
		},
		OFF: ModifyOttChannelStateState{
			value: "OFF",
		},
	}
}

func (c ModifyOttChannelStateState) Value() string {
	return c.value
}

func (c ModifyOttChannelStateState) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ModifyOttChannelStateState) UnmarshalJSON(b []byte) error {
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
