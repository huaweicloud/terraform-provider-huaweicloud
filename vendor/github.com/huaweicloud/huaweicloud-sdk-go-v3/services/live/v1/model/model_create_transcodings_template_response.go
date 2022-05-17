package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateTranscodingsTemplateResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateTranscodingsTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTranscodingsTemplateResponse struct{}"
	}

	return strings.Join([]string{"CreateTranscodingsTemplateResponse", string(data)}, " ")
}
