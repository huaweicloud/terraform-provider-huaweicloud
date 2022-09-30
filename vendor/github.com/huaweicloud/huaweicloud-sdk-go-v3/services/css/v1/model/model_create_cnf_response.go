package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateCnfResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateCnfResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCnfResponse struct{}"
	}

	return strings.Join([]string{"CreateCnfResponse", string(data)}, " ")
}
