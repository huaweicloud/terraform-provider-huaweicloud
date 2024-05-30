package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// DownloadErrorlogRequest Request Object
type DownloadErrorlogRequest struct {

	// 语言
	XLanguage *DownloadErrorlogRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o DownloadErrorlogRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DownloadErrorlogRequest struct{}"
	}

	return strings.Join([]string{"DownloadErrorlogRequest", string(data)}, " ")
}

type DownloadErrorlogRequestXLanguage struct {
	value string
}

type DownloadErrorlogRequestXLanguageEnum struct {
	ZH_CN DownloadErrorlogRequestXLanguage
	EN_US DownloadErrorlogRequestXLanguage
}

func GetDownloadErrorlogRequestXLanguageEnum() DownloadErrorlogRequestXLanguageEnum {
	return DownloadErrorlogRequestXLanguageEnum{
		ZH_CN: DownloadErrorlogRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: DownloadErrorlogRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c DownloadErrorlogRequestXLanguage) Value() string {
	return c.value
}

func (c DownloadErrorlogRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DownloadErrorlogRequestXLanguage) UnmarshalJSON(b []byte) error {
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
