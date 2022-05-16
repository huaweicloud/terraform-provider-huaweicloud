package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type IdentityToken struct {

	// token的ID。与请求头中的X-Auth-Token含义相同，待废弃。
	Id *string `json:"id,omitempty"`

	// AK/SK和securitytoken的有效期，时间单位为秒。取值范围：15min ~ 24h ，默认为15min。
	DurationSeconds *int32 `json:"duration_seconds,omitempty"`
}

func (o IdentityToken) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IdentityToken struct{}"
	}

	return strings.Join([]string{"IdentityToken", string(data)}, " ")
}
