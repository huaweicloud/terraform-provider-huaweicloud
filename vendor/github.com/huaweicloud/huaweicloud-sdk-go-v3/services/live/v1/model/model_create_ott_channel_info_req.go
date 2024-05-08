package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// CreateOttChannelInfoReq OTT频道创建消息体
type CreateOttChannelInfoReq struct {

	// 频道推流域名
	Domain string `json:"domain"`

	// 组名或应用名
	AppName string `json:"app_name"`

	// 频道ID。频道唯一标识，为必填项。频道ID不建议输入下划线“_”，否则会影响转码和截图任务
	Id string `json:"id"`

	// 频道名。可选配置
	Name *string `json:"name,omitempty"`

	// 频道状态 - ON：频道下发成功后，自动启动拉流、转码、录制等功能 - OFF：仅保存频道信息，不启动频道
	State CreateOttChannelInfoReqState `json:"state"`

	Input *InputStreamInfo `json:"input"`

	// 转码模板配置
	EncoderSettings *[]ModifyOttChannelEncoderSettingsEncoderSettings `json:"encoder_settings,omitempty"`

	RecordSettings *CreateOttChannelInfoReqRecordSettings `json:"record_settings"`

	// 频道出流信息
	Endpoints []EndpointItem `json:"endpoints"`
}

func (o CreateOttChannelInfoReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateOttChannelInfoReq struct{}"
	}

	return strings.Join([]string{"CreateOttChannelInfoReq", string(data)}, " ")
}

type CreateOttChannelInfoReqState struct {
	value string
}

type CreateOttChannelInfoReqStateEnum struct {
	ON  CreateOttChannelInfoReqState
	OFF CreateOttChannelInfoReqState
}

func GetCreateOttChannelInfoReqStateEnum() CreateOttChannelInfoReqStateEnum {
	return CreateOttChannelInfoReqStateEnum{
		ON: CreateOttChannelInfoReqState{
			value: "ON",
		},
		OFF: CreateOttChannelInfoReqState{
			value: "OFF",
		},
	}
}

func (c CreateOttChannelInfoReqState) Value() string {
	return c.value
}

func (c CreateOttChannelInfoReqState) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateOttChannelInfoReqState) UnmarshalJSON(b []byte) error {
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
