package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// VerifyDomainOwnerResponse Response Object
type VerifyDomainOwnerResponse struct {

	// 验证是否通过，true：通过，false：不通过。
	Result *bool `json:"result,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o VerifyDomainOwnerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VerifyDomainOwnerResponse struct{}"
	}

	return strings.Join([]string{"VerifyDomainOwnerResponse", string(data)}, " ")
}
