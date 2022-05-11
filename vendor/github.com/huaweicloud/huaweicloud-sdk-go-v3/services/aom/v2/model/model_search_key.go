package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 日志过滤条件集合，不同日志来源所需字段不同。
type SearchKey struct {

	// 应用名称。
	AppName *string `json:"appName,omitempty"`

	// CCE集群ID。
	ClusterId string `json:"clusterId"`

	// 日志所在虚拟机IP。
	HostIP *string `json:"hostIP,omitempty"`

	// CCE容器集群的命名空间。
	NameSpace *string `json:"nameSpace,omitempty"`

	// 日志文件名称。
	PathFile *string `json:"pathFile,omitempty"`

	// 容器实例名称。
	PodName *string `json:"podName,omitempty"`
}

func (o SearchKey) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SearchKey struct{}"
	}

	return strings.Join([]string{"SearchKey", string(data)}, " ")
}
