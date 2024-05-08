package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// VerifyDomainOwnerRequestBody 域名归属校验请求体
type VerifyDomainOwnerRequestBody struct {

	// 校验类型： - dns：DNS解析校验； - file：文件校验； - all：DNS与文件都会进行探测，默认为all。
	VerifyType *string `json:"verify_type,omitempty"`
}

func (o VerifyDomainOwnerRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VerifyDomainOwnerRequestBody struct{}"
	}

	return strings.Join([]string{"VerifyDomainOwnerRequestBody", string(data)}, " ")
}
