package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CertificatesResource struct {

	// 证书ID。
	Id *string `json:"id,omitempty"`

	// 证书名称。
	Name *string `json:"name,omitempty"`

	// SL证书的类型。分为服务器证书(server)、CA证书(client)。
	Type *string `json:"type,omitempty"`
}

func (o CertificatesResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CertificatesResource struct{}"
	}

	return strings.Join([]string{"CertificatesResource", string(data)}, " ")
}
