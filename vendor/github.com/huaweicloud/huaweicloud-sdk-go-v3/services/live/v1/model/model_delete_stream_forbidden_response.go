package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteStreamForbiddenResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteStreamForbiddenResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteStreamForbiddenResponse struct{}"
	}

	return strings.Join([]string{"DeleteStreamForbiddenResponse", string(data)}, " ")
}
