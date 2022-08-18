package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowDomainFullConfigResponse struct {
	Configs        *ConfigsGetBody `json:"configs,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ShowDomainFullConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainFullConfigResponse struct{}"
	}

	return strings.Join([]string{"ShowDomainFullConfigResponse", string(data)}, " ")
}
