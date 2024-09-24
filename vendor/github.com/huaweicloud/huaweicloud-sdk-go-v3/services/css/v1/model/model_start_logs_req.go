package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StartLogsReq struct {

	// 委托名称，委托给CSS，允许CSS调用您的其他云服务。
	Agency string `json:"agency"`

	// 日志在OBS桶中的备份路径。
	LogBasePath string `json:"logBasePath"`

	// 用于存储日志的OBS桶的桶名。
	LogBucket string `json:"logBucket"`

	// 保存日志的索引前缀。action等于real_time_log_collect时必选
	IndexPrefix *string `json:"index_prefix,omitempty"`

	// 日志保存时间。action等于real_time_log_collect时必选
	KeepDays *int32 `json:"keep_days,omitempty"`

	// 保存日志的目标集群。action等于real_time_log_collect时必选
	TargetClusterId *string `json:"target_cluster_id,omitempty"`
}

func (o StartLogsReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartLogsReq struct{}"
	}

	return strings.Join([]string{"StartLogsReq", string(data)}, " ")
}
