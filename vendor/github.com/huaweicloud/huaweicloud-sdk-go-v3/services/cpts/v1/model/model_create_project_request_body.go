package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateProjectRequestBody CreateProjectRequestBody
type CreateProjectRequestBody struct {

	// 名称
	Name string `json:"name"`

	// 描述
	Description *string `json:"description,omitempty"`
}

func (o CreateProjectRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateProjectRequestBody struct{}"
	}

	return strings.Join([]string{"CreateProjectRequestBody", string(data)}, " ")
}
