package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type RefreshTaskRequestBody struct {

	// 刷新的类型，其值可以为file：文件，或directory：目录，默认为file。
	Type *RefreshTaskRequestBodyType `json:"type,omitempty"`

	// 目录刷新方式，all：刷新目录下全部资源；detect_modify_refresh：刷新目录下已变更的资源，默认值为all。
	Mode *RefreshTaskRequestBodyMode `json:"mode,omitempty"`

	// 是否对url中的中文字符进行编码后刷新，false代表不开启，true代表开启，开启后仅刷新转码后的URL。
	ZhUrlEncode *bool `json:"zh_url_encode,omitempty"`

	// 需要刷新的URL必须带有“http://”或“https://”，多个URL用逗号分隔（\"url1\", \"url2\"），单个url的长度限制为4096字符，单次最多输入1000个url，如果输入的是目录，支持100个目录刷新。   > - 如果您需要刷新的URL中有中文，请同时刷新中文URL（输入中文URL且不开启zh_url_encode）和转码后的URL（输入中文URL且开启zh_url_encode）。   > - 如果您的URL中带有空格，请自行转码后输入，且不要开启URL Encode。
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
