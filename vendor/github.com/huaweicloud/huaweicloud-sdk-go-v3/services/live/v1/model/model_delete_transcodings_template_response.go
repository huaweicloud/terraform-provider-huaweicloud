package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteTranscodingsTemplateResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteTranscodingsTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTranscodingsTemplateResponse struct{}"
	}

	return strings.Join([]string{"DeleteTranscodingsTemplateResponse", string(data)}, " ")
}
