package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateInstanceUserResponse Response Object
type UpdateInstanceUserResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateInstanceUserResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceUserResponse struct{}"
	}

	return strings.Join([]string{"UpdateInstanceUserResponse", string(data)}, " ")
}
