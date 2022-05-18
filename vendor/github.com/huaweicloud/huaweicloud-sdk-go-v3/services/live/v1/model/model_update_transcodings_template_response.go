package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateTranscodingsTemplateResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateTranscodingsTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTranscodingsTemplateResponse struct{}"
	}

	return strings.Join([]string{"UpdateTranscodingsTemplateResponse", string(data)}, " ")
}
