package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ReleaseReqBodyParams 模板实例参数
type ReleaseReqBodyParams struct {

	// 开启后，仅验证模板参数，不进行安装
	DryRun *bool `json:"dry_run,omitempty"`

	// 实例名称模板
	NameTemplate *string `json:"name_template,omitempty"`

	// 安装时是否禁用hooks
	NoHooks *bool `json:"no_hooks,omitempty"`

	// 是否替换同名实例
	Replace *bool `json:"replace,omitempty"`

	// 是否重建实例
	Recreate *bool `json:"recreate,omitempty"`

	// 更新时是否重置values
	ResetValues *bool `json:"reset_values,omitempty"`

	// 回滚实例的版本
	ReleaseVersion *int32 `json:"release_version,omitempty"`

	// 更新或者删除时启用hooks
	IncludeHooks *bool `json:"include_hooks,omitempty"`
}

func (o ReleaseReqBodyParams) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReleaseReqBodyParams struct{}"
	}

	return strings.Join([]string{"ReleaseReqBodyParams", string(data)}, " ")
}
