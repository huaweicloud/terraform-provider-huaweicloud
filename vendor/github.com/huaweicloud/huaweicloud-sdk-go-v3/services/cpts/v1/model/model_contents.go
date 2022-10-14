package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Contents struct {

	// 事务id，若不为0表示此卡片为事务；为0表示非事务
	ContentId *int32 `json:"content_id,omitempty"`

	// content
	Content *[]Content `json:"content,omitempty"`

	// 排序索引标识
	Index *int32 `json:"index,omitempty"`

	// selected_temp_name
	SelectedTempName *string `json:"selected_temp_name,omitempty"`

	// 数据（循环、条件控制器作用的数据）
	Data *interface{} `json:"data,omitempty"`

	// 类型，0:默认请求；1:数据指令；201:循环指令； 202:条件指令；301:集合点
	DataType *int32 `json:"data_type,omitempty"`

	// 若类型为202:条件指令，该字段为条件配置
	Conditions *interface{} `json:"conditions,omitempty"`

	// 是否禁用
	IsDisabled *bool `json:"is_disabled,omitempty"`
}

func (o Contents) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Contents struct{}"
	}

	return strings.Join([]string{"Contents", string(data)}, " ")
}
