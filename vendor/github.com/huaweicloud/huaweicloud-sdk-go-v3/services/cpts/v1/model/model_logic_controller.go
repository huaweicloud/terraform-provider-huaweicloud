package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type LogicController struct {
	// for_loop_params

	ForLoopParams *string `json:"for_loop_params,omitempty"`
	// condition

	Condition *string `json:"condition,omitempty"`
}

func (o LogicController) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LogicController struct{}"
	}

	return strings.Join([]string{"LogicController", string(data)}, " ")
}
