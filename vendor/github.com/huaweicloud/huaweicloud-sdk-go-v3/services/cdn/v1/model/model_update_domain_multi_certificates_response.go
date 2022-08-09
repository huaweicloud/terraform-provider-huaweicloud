package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateDomainMultiCertificatesResponse struct {
	Https          *UpdateDomainMultiCertificatesResponseBodyContent `json:"https,omitempty"`
	HttpStatusCode int                                               `json:"-"`
}

func (o UpdateDomainMultiCertificatesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainMultiCertificatesResponse struct{}"
	}

	return strings.Join([]string{"UpdateDomainMultiCertificatesResponse", string(data)}, " ")
}
