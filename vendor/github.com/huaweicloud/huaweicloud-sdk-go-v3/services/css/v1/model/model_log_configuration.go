package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type LogConfiguration struct {

	// 日志备份ID，通过系统UUID生成。
	Id *string `json:"id,omitempty"`

	// 集群ID。
	ClusterId *string `json:"clusterId,omitempty"`

	// 用于存储日志的OBS桶的桶名。
	ObsBucket *string `json:"obsBucket,omitempty"`

	// 委托名称，委托给CSS，允许CSS调用您的其他云服务。
	Agency *string `json:"agency,omitempty"`

	// 更新时间。格式为：Unix时间戳格式。
	UpdateAt *int64 `json:"updateAt,omitempty"`

	// 日志在OBS桶中的备份路径。
	BasePath *string `json:"basePath,omitempty"`

	// 自动备份开关。 - true: 自动备份开启。 - false: 自动备份关闭。
	AutoEnable *bool `json:"autoEnable,omitempty"`

	// 自动备份日志开始时间。当autoEnable为false时该字段为null。格式为：格林威治标准时间。
	Period *string `json:"period,omitempty"`

	// 日志开关。 - true: 日志开启。 - false: 日志关闭。
	LogSwitch *bool `json:"logSwitch,omitempty"`
}

func (o LogConfiguration) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LogConfiguration struct{}"
	}

	return strings.Join([]string{"LogConfiguration", string(data)}, " ")
}
