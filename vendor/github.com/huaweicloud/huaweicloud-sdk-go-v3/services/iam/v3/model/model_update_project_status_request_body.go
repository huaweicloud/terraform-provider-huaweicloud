package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateProjectStatusRequestBody struct {
	Project *UpdateProjectOption `json:"project"`
}

func (o UpdateProjectStatusRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProjectStatusRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateProjectStatusRequestBody", string(data)}, " ")
}
