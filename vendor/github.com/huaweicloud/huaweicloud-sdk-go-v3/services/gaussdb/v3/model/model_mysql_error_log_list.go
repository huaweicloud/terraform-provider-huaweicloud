package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlErrorLogList struct {
	// 节点ID。

	NodeId *string `json:"node_id,omitempty"`
	// 日期时间UTC时间。

	Time *string `json:"time,omitempty"`
	// 日志级别。

	Level *string `json:"level,omitempty"`
	// 错误日志内容。

	Content *string `json:"content,omitempty"`
}

func (o MysqlErrorLogList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlErrorLogList struct{}"
	}

	return strings.Join([]string{"MysqlErrorLogList", string(data)}, " ")
}
