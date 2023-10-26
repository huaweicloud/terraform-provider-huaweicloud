package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListElbCertsResponse Response Object
type ListElbCertsResponse struct {
	Certificates   *CertificatesResource `json:"certificates,omitempty"`
	HttpStatusCode int                   `json:"-"`
}

func (o ListElbCertsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListElbCertsResponse struct{}"
	}

	return strings.Join([]string{"ListElbCertsResponse", string(data)}, " ")
}
