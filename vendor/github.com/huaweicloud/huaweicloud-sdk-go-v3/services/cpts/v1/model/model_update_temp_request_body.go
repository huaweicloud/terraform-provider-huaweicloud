package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTempRequestBody
type UpdateTempRequestBody struct {
	// id

	Id int32 `json:"id"`
	// project_id

	ProjectId int32 `json:"project_id"`
	// name

	Name string `json:"name"`
	// temp_type

	TempType *int32 `json:"temp_type,omitempty"`
	// description

	Description *string `json:"description,omitempty"`
	// for_loop_params

	ForLoopParams *[]interface{} `json:"for_loop_params,omitempty"`
	// enable_pre

	EnablePre *bool `json:"enable_pre,omitempty"`
	// contents

	Contents *[]TempContentInfo `json:"contents,omitempty"`
}

func (o UpdateTempRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTempRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateTempRequestBody", string(data)}, " ")
}
