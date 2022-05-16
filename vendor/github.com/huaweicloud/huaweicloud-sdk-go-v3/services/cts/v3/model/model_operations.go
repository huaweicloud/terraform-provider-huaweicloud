package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 操作事件列表，目前最多支持对100服务的1000个关键操作进行配置。
type Operations struct {

	// 标识云服务类型。必须为已对接CTS的云服务的英文缩写，且服务类型一般为大写字母。 已对接的云服务列表参见《云审计服务用户指南》“支持的服务”章节。
	ServiceType string `json:"service_type"`

	// 标识资源类型。
	ResourceType string `json:"resource_type"`

	// 标识事件名称。
	TraceNames []string `json:"trace_names"`
}

func (o Operations) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Operations struct{}"
	}

	return strings.Join([]string{"Operations", string(data)}, " ")
}
