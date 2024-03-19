package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type LogObject struct {

	// 域名名称。
	DomainName *string `json:"domain_name,omitempty"`

	// 查询起始时间，相对于UTC 1970-01-01到当前时间相隔的毫秒数。
	StartTime *int64 `json:"start_time,omitempty"`

	// 查询结束时间，相对于UTC 1970-01-01到当前时间相隔的毫秒数。
	EndTime *int64 `json:"end_time,omitempty"`

	// 日志文件名字。
	Name *string `json:"name,omitempty"`

	// 文件大小(Byte)。
	Size *int64 `json:"size,omitempty"`

	// 下载链接。
	Link *string `json:"link,omitempty"`
}

func (o LogObject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LogObject struct{}"
	}

	return strings.Join([]string{"LogObject", string(data)}, " ")
}
