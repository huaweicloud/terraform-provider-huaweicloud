package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateProjectRequestBody
type UpdateProjectRequestBody struct {
	// id

	Id int32 `json:"id"`
	// name

	Name string `json:"name"`
	// description

	Description *string `json:"description,omitempty"`
	// variables_no_file

	VariablesNoFile *[]string `json:"variables_no_file,omitempty"`
	// source

	Source *int32 `json:"source,omitempty"`
	// external_params

	ExternalParams *interface{} `json:"external_params,omitempty"`
}

func (o UpdateProjectRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProjectRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateProjectRequestBody", string(data)}, " ")
}
