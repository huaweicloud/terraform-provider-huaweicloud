package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// id token信息
type GetIdTokenIdTokenBody struct {

	// id_token的值
	Id string `json:"id"`
}

func (o GetIdTokenIdTokenBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetIdTokenIdTokenBody struct{}"
	}

	return strings.Join([]string{"GetIdTokenIdTokenBody", string(data)}, " ")
}
