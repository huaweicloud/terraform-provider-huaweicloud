package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ConfigurationMetadata metadata字段数据结构说明
type ConfigurationMetadata struct {

	// Configuration名称
	Name string `json:"name"`

	// Configuration标签，key/value对格式。  - Key：必须以字母或数字开头，可以包含字母、数字、连字符、下划线和点，最长63个字符；另外可以使用DNS子域作为前缀，例如example.com/my-key，DNS子域最长253个字符。 - Value：可以为空或者非空字符串，非空字符串必须以字符或数字开头，可以包含字母、数字、连字符、下划线和点，最长63个字符。  示例：\"foo\": \"bar\"
	Labels map[string]string `json:"labels,omitempty"`
}

func (o ConfigurationMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfigurationMetadata struct{}"
	}

	return strings.Join([]string{"ConfigurationMetadata", string(data)}, " ")
}
