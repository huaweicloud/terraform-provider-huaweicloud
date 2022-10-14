package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type StreamTranscodingTemplate struct {

	// 推流域名
	Domain string `json:"domain"`

	// 应用名称。 默认为“live”，若您需要自定义应用名称，请先提交工单申请。
	AppName string `json:"app_name"`

	// 转码流触发模式。 - play：拉流触发转码。 - publish：推流触发转码。 默认为play
	TransType *StreamTranscodingTemplateTransType `json:"trans_type,omitempty"`

	// 视频质量信息
	QualityInfo []QualityInfo `json:"quality_info"`
}

func (o StreamTranscodingTemplate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StreamTranscodingTemplate struct{}"
	}

	return strings.Join([]string{"StreamTranscodingTemplate", string(data)}, " ")
}

type StreamTranscodingTemplateTransType struct {
	value string
}

type StreamTranscodingTemplateTransTypeEnum struct {
	PLAY    StreamTranscodingTemplateTransType
	PUBLISH StreamTranscodingTemplateTransType
}

func GetStreamTranscodingTemplateTransTypeEnum() StreamTranscodingTemplateTransTypeEnum {
	return StreamTranscodingTemplateTransTypeEnum{
		PLAY: StreamTranscodingTemplateTransType{
			value: "play",
		},
		PUBLISH: StreamTranscodingTemplateTransType{
			value: "publish",
		},
	}
}

func (c StreamTranscodingTemplateTransType) Value() string {
	return c.value
}

func (c StreamTranscodingTemplateTransType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *StreamTranscodingTemplateTransType) UnmarshalJSON(b []byte) error {
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
