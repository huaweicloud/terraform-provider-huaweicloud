package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateLogSettingReq struct {

	// 委托名称，委托给CSS，允许CSS调用您的其他云服务。
	Agency string `json:"agency"`

	// 日志在OBS桶中的备份路径。
	LogBasePath string `json:"logBasePath"`

	// 用于存储日志的OBS桶的桶名。
	LogBucket string `json:"logBucket"`
}

func (o UpdateLogSettingReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateLogSettingReq struct{}"
	}

	return strings.Join([]string{"UpdateLogSettingReq", string(data)}, " ")
}
