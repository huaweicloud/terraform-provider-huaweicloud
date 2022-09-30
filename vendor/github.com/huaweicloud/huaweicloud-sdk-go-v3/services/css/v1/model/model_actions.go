package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Actions struct {

	// 操作记录id。
	Id *string `json:"id,omitempty"`

	// 操作类型。
	ActionType *string `json:"actionType,omitempty"`

	// 配置文件内容。
	ConfContent *string `json:"confContent,omitempty"`

	// 操作状态。
	Status *string `json:"status,omitempty"`

	// 更新时间，格式为ISO8601：CCYY-MM-DDThh:mm:ss。
	UpdateAt *string `json:"updateAt,omitempty"`

	// 错误信息。当操作状态为success时该字段为null。
	ErrorMsg *string `json:"errorMsg,omitempty"`

	// 内容。
	Message *string `json:"message,omitempty"`
}

func (o Actions) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Actions struct{}"
	}

	return strings.Join([]string{"Actions", string(data)}, " ")
}
