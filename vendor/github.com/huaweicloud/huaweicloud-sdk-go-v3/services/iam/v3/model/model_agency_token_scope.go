package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type AgencyTokenScope struct {
	Domain *AgencyTokenScopeDomain `json:"domain,omitempty"`

	Project *AgencyTokenScopeProject `json:"project,omitempty"`
}

func (o AgencyTokenScope) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AgencyTokenScope struct{}"
	}

	return strings.Join([]string{"AgencyTokenScope", string(data)}, " ")
}
