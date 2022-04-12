package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateTempResponse struct {
	// code

	Code *string `json:"code,omitempty"`
	// tempId

	TempId *int32 `json:"tempId,omitempty"`
	// message

	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateTempResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTempResponse struct{}"
	}

	return strings.Join([]string{"CreateTempResponse", string(data)}, " ")
}
