package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type RefreshTaskRequestBody struct {

	// 刷新的类型，其值可以为file 或directory，默认为file
	Type *RefreshTaskRequestBodyType `json:"type,omitempty"`

	// 目录刷新方式，all：刷新目录下全部资源；detect_modify_refresh：刷新目录下已变更的资源，默认值为all。
	Mode *RefreshTaskRequestBodyMode `json:"mode,omitempty"`

	// 输入URL必须带有“http://”或“https://”，多个URL用逗号分隔，单个url的长度限制为4096字符，单次最多输入1000个url。 >   如果您需要刷新的URL中有中文，请同时刷新中文URL和转码后的URL。
	Urls []string `json:"urls"`
}

func (o RefreshTaskRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RefreshTaskRequestBody struct{}"
	}

	return strings.Join([]string{"RefreshTaskRequestBody", string(data)}, " ")
}

type RefreshTaskRequestBodyType struct {
	value string
}

type RefreshTaskRequestBodyTypeEnum struct {
	FILE      RefreshTaskRequestBodyType
	DIRECTORY RefreshTaskRequestBodyType
}

func GetRefreshTaskRequestBodyTypeEnum() RefreshTaskRequestBodyTypeEnum {
	return RefreshTaskRequestBodyTypeEnum{
		FILE: RefreshTaskRequestBodyType{
			value: "file",
		},
		DIRECTORY: RefreshTaskRequestBodyType{
			value: "directory",
		},
	}
}

func (c RefreshTaskRequestBodyType) Value() string {
	return c.value
}

func (c RefreshTaskRequestBodyType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RefreshTaskRequestBodyType) UnmarshalJSON(b []byte) error {
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

type RefreshTaskRequestBodyMode struct {
	value string
}

type RefreshTaskRequestBodyModeEnum struct {
	ALL                   RefreshTaskRequestBodyMode
	DETECT_MODIFY_REFRESH RefreshTaskRequestBodyMode
}

func GetRefreshTaskRequestBodyModeEnum() RefreshTaskRequestBodyModeEnum {
	return RefreshTaskRequestBodyModeEnum{
		ALL: RefreshTaskRequestBodyMode{
			value: "all",
		},
		DETECT_MODIFY_REFRESH: RefreshTaskRequestBodyMode{
			value: "detect_modify_refresh",
		},
	}
}

func (c RefreshTaskRequestBodyMode) Value() string {
	return c.value
}

func (c RefreshTaskRequestBodyMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RefreshTaskRequestBodyMode) UnmarshalJSON(b []byte) error {
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
