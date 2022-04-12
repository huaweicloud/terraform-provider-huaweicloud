package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TempInfo struct {
	// id

	Id *int32 `json:"id,omitempty"`
	// project_id

	ProjectId *int32 `json:"project_id,omitempty"`
	// name

	Name *string `json:"name,omitempty"`
	// description

	Description *string `json:"description,omitempty"`
	// variables

	Variables *string `json:"variables,omitempty"`
	// contents

	Contents *[]interface{} `json:"contents,omitempty"`
	// temp_type

	TempType *int32 `json:"temp_type,omitempty"`
	// for_loop_params

	ForLoopParams *[]interface{} `json:"for_loop_params,omitempty"`

	LogicController *LogicController `json:"logic_controller,omitempty"`
	// enable_pre

	EnablePre *bool `json:"enable_pre,omitempty"`
}

func (o TempInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TempInfo struct{}"
	}

	return strings.Join([]string{"TempInfo", string(data)}, " ")
}
