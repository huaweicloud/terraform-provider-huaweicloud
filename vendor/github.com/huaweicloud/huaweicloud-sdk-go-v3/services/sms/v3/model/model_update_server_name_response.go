package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateServerNameResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateServerNameResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateServerNameResponse struct{}"
	}

	return strings.Join([]string{"UpdateServerNameResponse", string(data)}, " ")
}
