package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDomainHttpsCertResponse Response Object
type DeleteDomainHttpsCertResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteDomainHttpsCertResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDomainHttpsCertResponse struct{}"
	}

	return strings.Join([]string{"DeleteDomainHttpsCertResponse", string(data)}, " ")
}
