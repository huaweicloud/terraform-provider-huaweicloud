package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type AgencyAuth struct {
	Identity *AgencyAuthIdentity `json:"identity"`
}

func (o AgencyAuth) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AgencyAuth struct{}"
	}

	return strings.Join([]string{"AgencyAuth", string(data)}, " ")
}
