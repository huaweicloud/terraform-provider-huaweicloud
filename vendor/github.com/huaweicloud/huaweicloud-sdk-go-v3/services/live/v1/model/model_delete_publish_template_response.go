package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePublishTemplateResponse Response Object
type DeletePublishTemplateResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeletePublishTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePublishTemplateResponse struct{}"
	}

	return strings.Join([]string{"DeletePublishTemplateResponse", string(data)}, " ")
}
