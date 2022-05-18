package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateStreamForbiddenResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateStreamForbiddenResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateStreamForbiddenResponse struct{}"
	}

	return strings.Join([]string{"UpdateStreamForbiddenResponse", string(data)}, " ")
}
