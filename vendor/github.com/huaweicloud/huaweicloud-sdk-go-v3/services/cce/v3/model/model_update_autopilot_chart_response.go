package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAutopilotChartResponse Response Object
type UpdateAutopilotChartResponse struct {

	// 模板ID
	Id *string `json:"id,omitempty"`

	// 模板名称
	Name *string `json:"name,omitempty"`

	// 模板值
	Values *string `json:"values,omitempty"`

	// 模板翻译资源
	Translate *string `json:"translate,omitempty"`

	// 模板介绍
	Instruction *string `json:"instruction,omitempty"`

	// 模板版本
	Version *string `json:"version,omitempty"`

	// 模板描述
	Description *string `json:"description,omitempty"`

	// 模板的来源
	Source *string `json:"source,omitempty"`

	// 模板的图标链接
	IconUrl *string `json:"icon_url,omitempty"`

	// 是否公开模板
	Public *bool `json:"public,omitempty"`

	// 模板的链接
	ChartUrl *string `json:"chart_url,omitempty"`

	// 创建时间
	CreateAt *string `json:"create_at,omitempty"`

	// 更新时间
	UpdateAt       *string `json:"update_at,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateAutopilotChartResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAutopilotChartResponse struct{}"
	}

	return strings.Join([]string{"UpdateAutopilotChartResponse", string(data)}, " ")
}
