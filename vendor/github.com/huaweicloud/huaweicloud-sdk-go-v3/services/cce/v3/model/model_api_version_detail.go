package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ApiVersionDetail API版本的详细信息
type ApiVersionDetail struct {

	// API版本ID。例如v3。
	Id string `json:"id"`

	// API版本的URL链接信息。
	Links []ApiVersionLink `json:"links"`

	// 如果API的这个版本支持微版本，则支持最小的微版本。如果不支持微版本，这将是空字符串。
	MinVersion string `json:"min_version"`

	// API版本的状态。 可以是： - CURRENT这是使用的API的首选版本； - SUPPORTED：这是一个较老的，但仍然支持的API版本； - DEPRECATED：一个被废弃的API版本，该版本将被删除
	Status ApiVersionDetailStatus `json:"status"`

	// API发布时间（UTC格式）。例如API版本为v3时，值为'2018-09-15 00:00:00Z'。
	Updated string `json:"updated"`

	// 如果API的这个版本支持微版本，则支持最大的微版本。如果不支持微版本，这将是空字符串。
	Version string `json:"version"`
}

func (o ApiVersionDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ApiVersionDetail struct{}"
	}

	return strings.Join([]string{"ApiVersionDetail", string(data)}, " ")
}

type ApiVersionDetailStatus struct {
	value string
}

type ApiVersionDetailStatusEnum struct {
	CURRENT    ApiVersionDetailStatus
	SUPPORTED  ApiVersionDetailStatus
	DEPRECATED ApiVersionDetailStatus
}

func GetApiVersionDetailStatusEnum() ApiVersionDetailStatusEnum {
	return ApiVersionDetailStatusEnum{
		CURRENT: ApiVersionDetailStatus{
			value: "CURRENT",
		},
		SUPPORTED: ApiVersionDetailStatus{
			value: "SUPPORTED",
		},
		DEPRECATED: ApiVersionDetailStatus{
			value: "DEPRECATED",
		},
	}
}

func (c ApiVersionDetailStatus) Value() string {
	return c.value
}

func (c ApiVersionDetailStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ApiVersionDetailStatus) UnmarshalJSON(b []byte) error {
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
