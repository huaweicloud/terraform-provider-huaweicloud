package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UnscopedTokenInfoRoles struct {

	// role id。
	Id *string `json:"id,omitempty"`

	// name id。
	Name *string `json:"name,omitempty"`
}

func (o UnscopedTokenInfoRoles) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnscopedTokenInfoRoles struct{}"
	}

	return strings.Join([]string{"UnscopedTokenInfoRoles", string(data)}, " ")
}
