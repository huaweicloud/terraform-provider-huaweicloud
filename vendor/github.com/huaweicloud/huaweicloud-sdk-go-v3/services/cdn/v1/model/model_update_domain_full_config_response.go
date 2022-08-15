package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateDomainFullConfigResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateDomainFullConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainFullConfigResponse struct{}"
	}

	return strings.Join([]string{"UpdateDomainFullConfigResponse", string(data)}, " ")
}
