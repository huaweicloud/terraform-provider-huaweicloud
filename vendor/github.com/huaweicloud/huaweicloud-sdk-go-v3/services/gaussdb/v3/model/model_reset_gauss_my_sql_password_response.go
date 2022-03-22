package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ResetGaussMySqlPasswordResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ResetGaussMySqlPasswordResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetGaussMySqlPasswordResponse struct{}"
	}

	return strings.Join([]string{"ResetGaussMySqlPasswordResponse", string(data)}, " ")
}
