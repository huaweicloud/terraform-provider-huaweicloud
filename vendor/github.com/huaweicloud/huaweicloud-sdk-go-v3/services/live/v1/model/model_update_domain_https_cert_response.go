package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDomainHttpsCertResponse Response Object
type UpdateDomainHttpsCertResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateDomainHttpsCertResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainHttpsCertResponse struct{}"
	}

	return strings.Join([]string{"UpdateDomainHttpsCertResponse", string(data)}, " ")
}
