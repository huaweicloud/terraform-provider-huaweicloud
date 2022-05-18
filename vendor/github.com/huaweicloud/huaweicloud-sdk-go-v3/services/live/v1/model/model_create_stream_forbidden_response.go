package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateStreamForbiddenResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateStreamForbiddenResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateStreamForbiddenResponse struct{}"
	}

	return strings.Join([]string{"CreateStreamForbiddenResponse", string(data)}, " ")
}
