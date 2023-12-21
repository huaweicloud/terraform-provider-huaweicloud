package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListPublishTemplateResponse Response Object
type ListPublishTemplateResponse struct {

	// 回调地址
	Url *string `json:"url,omitempty"`

	// 鉴权密钥
	AuthSignKey *string `json:"auth_sign_key,omitempty"`

	// 接收回调通知服务器所在区域。 包含如下取值： - mainland_china：中国大陆区域。 - outside_mainland_china：中国大陆以外区域。
	CallBackArea   *ListPublishTemplateResponseCallBackArea `json:"call_back_area,omitempty"`
	HttpStatusCode int                                      `json:"-"`
}

func (o ListPublishTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPublishTemplateResponse struct{}"
	}

	return strings.Join([]string{"ListPublishTemplateResponse", string(data)}, " ")
}

type ListPublishTemplateResponseCallBackArea struct {
	value string
}

type ListPublishTemplateResponseCallBackAreaEnum struct {
	MAINLAND_CHINA         ListPublishTemplateResponseCallBackArea
	OUTSIDE_MAINLAND_CHINA ListPublishTemplateResponseCallBackArea
}

func GetListPublishTemplateResponseCallBackAreaEnum() ListPublishTemplateResponseCallBackAreaEnum {
	return ListPublishTemplateResponseCallBackAreaEnum{
		MAINLAND_CHINA: ListPublishTemplateResponseCallBackArea{
			value: "mainland_china",
		},
		OUTSIDE_MAINLAND_CHINA: ListPublishTemplateResponseCallBackArea{
			value: "outside_mainland_china",
		},
	}
}

func (c ListPublishTemplateResponseCallBackArea) Value() string {
	return c.value
}

func (c ListPublishTemplateResponseCallBackArea) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListPublishTemplateResponseCallBackArea) UnmarshalJSON(b []byte) error {
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
