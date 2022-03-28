package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// auth信息
type GetIdTokenAuthParams struct {
	IdToken *GetIdTokenIdTokenBody `json:"id_token"`

	Scope *GetIdTokenIdScopeBody `json:"scope,omitempty"`
}

func (o GetIdTokenAuthParams) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetIdTokenAuthParams struct{}"
	}

	return strings.Join([]string{"GetIdTokenAuthParams", string(data)}, " ")
}
