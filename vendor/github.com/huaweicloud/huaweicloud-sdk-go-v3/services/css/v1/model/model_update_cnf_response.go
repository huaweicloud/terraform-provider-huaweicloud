package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateCnfResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateCnfResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCnfResponse struct{}"
	}

	return strings.Join([]string{"UpdateCnfResponse", string(data)}, " ")
}
