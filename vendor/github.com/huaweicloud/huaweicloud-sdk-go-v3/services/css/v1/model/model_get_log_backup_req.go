package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type GetLogBackupReq struct {

	// 节点名称。通过[查询集群详情](ShowClusterDetail.xml)获取instances中的name属性。
	InstanceName string `json:"instanceName"`

	// 日志级别。可查询的日志级别为：INFO，ERROR，DEBUG，WARN。
	Level string `json:"level"`

	// 日志类型。可查询的日志类型为：deprecation，indexingSlow，searchSlow， instance。
	LogType string `json:"logType"`
}

func (o GetLogBackupReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetLogBackupReq struct{}"
	}

	return strings.Join([]string{"GetLogBackupReq", string(data)}, " ")
}
