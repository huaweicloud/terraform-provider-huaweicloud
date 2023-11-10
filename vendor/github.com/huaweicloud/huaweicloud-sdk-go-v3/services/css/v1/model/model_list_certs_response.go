package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListCertsResponse Response Object
type ListCertsResponse struct {
	DefaultCerts *DefaultCertsResource `json:"defaultCerts,omitempty"`

	CustomCerts    *CustomCertsResource `json:"customCerts,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o ListCertsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListCertsResponse struct{}"
	}

	return strings.Join([]string{"ListCertsResponse", string(data)}, " ")
}
