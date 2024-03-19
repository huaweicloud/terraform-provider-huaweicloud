package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowVerifyDomainOwnerInfoResponse Response Object
type ShowVerifyDomainOwnerInfoResponse struct {

	// DNS探测类型
	DnsVerifyType *string `json:"dns_verify_type,omitempty"`

	// DNS记录名称
	DnsVerifyName *string `json:"dns_verify_name,omitempty"`

	// 文件探测地址
	FileVerifyUrl *string `json:"file_verify_url,omitempty"`

	// 域名
	DomainName *string `json:"domain_name,omitempty"`

	// 探测域名
	VerifyDomainName *string `json:"verify_domain_name,omitempty"`

	// 探测文件名
	FileVerifyFilename *string `json:"file_verify_filename,omitempty"`

	// 探测内容，DNS值或者文件内容，时间加uuid
	VerifyContent *string `json:"verify_content,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowVerifyDomainOwnerInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowVerifyDomainOwnerInfoResponse struct{}"
	}

	return strings.Join([]string{"ShowVerifyDomainOwnerInfoResponse", string(data)}, " ")
}
