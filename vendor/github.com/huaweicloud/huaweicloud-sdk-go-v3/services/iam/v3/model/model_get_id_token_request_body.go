package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 请求体
type GetIdTokenRequestBody struct {
	Auth *GetIdTokenAuthParams `json:"auth"`
}

func (o GetIdTokenRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetIdTokenRequestBody struct{}"
	}

	return strings.Join([]string{"GetIdTokenRequestBody", string(data)}, " ")
}
