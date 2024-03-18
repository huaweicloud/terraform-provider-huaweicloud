package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type LogicController struct {

	// 旧版本逻辑控制器字段，当前已未使用
	ForLoopParams *string `json:"for_loop_params,omitempty"`

	// 逻辑控制器条件
	Condition *string `json:"condition,omitempty"`
}

func (o LogicController) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LogicController struct{}"
	}

	return strings.Join([]string{"LogicController", string(data)}, " ")
}
