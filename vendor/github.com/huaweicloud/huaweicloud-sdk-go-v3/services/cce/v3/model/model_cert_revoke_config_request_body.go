package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CertRevokeConfigRequestBody struct {

	// 用户ID
	UserId *string `json:"userId,omitempty"`

	// 委托用户ID
	AgencyId *string `json:"agencyId,omitempty"`
}

func (o CertRevokeConfigRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CertRevokeConfigRequestBody struct{}"
	}

	return strings.Join([]string{"CertRevokeConfigRequestBody", string(data)}, " ")
}
