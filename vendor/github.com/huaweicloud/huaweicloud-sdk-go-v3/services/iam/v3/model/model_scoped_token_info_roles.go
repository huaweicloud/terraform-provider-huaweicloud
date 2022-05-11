package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ScopedTokenInfoRoles struct {

	// role id
	Id *string `json:"id,omitempty"`

	// name id
	Name *string `json:"name,omitempty"`
}

func (o ScopedTokenInfoRoles) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScopedTokenInfoRoles struct{}"
	}

	return strings.Join([]string{"ScopedTokenInfoRoles", string(data)}, " ")
}
