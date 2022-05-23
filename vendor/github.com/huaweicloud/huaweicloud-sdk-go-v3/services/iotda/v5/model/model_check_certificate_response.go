package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CheckCertificateResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CheckCertificateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckCertificateResponse struct{}"
	}

	return strings.Join([]string{"CheckCertificateResponse", string(data)}, " ")
}
