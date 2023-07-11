package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResetUserPasswrodResponse Response Object
type ResetUserPasswrodResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ResetUserPasswrodResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetUserPasswrodResponse struct{}"
	}

	return strings.Join([]string{"ResetUserPasswrodResponse", string(data)}, " ")
}
