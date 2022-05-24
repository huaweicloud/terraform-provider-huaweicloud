package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateTemplateGroupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateTemplateGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTemplateGroupResponse struct{}"
	}

	return strings.Join([]string{"UpdateTemplateGroupResponse", string(data)}, " ")
}
