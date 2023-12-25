package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ApiVersionLink API版本的URL链接信息。
type ApiVersionLink struct {

	// API版本信息的链接。
	Href string `json:"href"`

	// 链接属性。self：自助链接包含版本链接的资源。立即链接后使用这些链接。
	Rel ApiVersionLinkRel `json:"rel"`
}

func (o ApiVersionLink) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ApiVersionLink struct{}"
	}

	return strings.Join([]string{"ApiVersionLink", string(data)}, " ")
}

type ApiVersionLinkRel struct {
	value string
}

type ApiVersionLinkRelEnum struct {
	SELF ApiVersionLinkRel
}

func GetApiVersionLinkRelEnum() ApiVersionLinkRelEnum {
	return ApiVersionLinkRelEnum{
		SELF: ApiVersionLinkRel{
			value: "self",
		},
	}
}

func (c ApiVersionLinkRel) Value() string {
	return c.value
}

func (c ApiVersionLinkRel) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ApiVersionLinkRel) UnmarshalJSON(b []byte) error {
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
