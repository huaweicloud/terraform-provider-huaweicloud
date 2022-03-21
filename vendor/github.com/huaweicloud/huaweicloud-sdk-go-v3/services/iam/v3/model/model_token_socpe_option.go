package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type TokenSocpeOption struct {
	Domain *ScopeDomainOption `json:"domain,omitempty"`

	Project *ScopeProjectOption `json:"project,omitempty"`
}

func (o TokenSocpeOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TokenSocpeOption struct{}"
	}

	return strings.Join([]string{"TokenSocpeOption", string(data)}, " ")
}
