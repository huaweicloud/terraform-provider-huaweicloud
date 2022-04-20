package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateCaseRequestBody
type UpdateCaseRequestBody struct {
	// contents

	Contents *[]CaseInfo `json:"contents,omitempty"`
	// for_loop_params

	ForLoopParams *[]interface{} `json:"for_loop_params,omitempty"`
}

func (o UpdateCaseRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCaseRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateCaseRequestBody", string(data)}, " ")
}
