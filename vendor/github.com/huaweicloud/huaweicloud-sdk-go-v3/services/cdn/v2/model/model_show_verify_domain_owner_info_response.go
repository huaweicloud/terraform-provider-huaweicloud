package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowVerifyDomainOwnerInfoResponse Response Object
type ShowVerifyDomainOwnerInfoResponse struct {

	// DNS解析类型。
	DnsVerifyType *string `json:"dns_verify_type,omitempty"`

	// DNS解析主机记录名称。
	DnsVerifyName *string `json:"dns_verify_name,omitempty"`

	// 文件校验URL地址。
	FileVerifyUrl *string `json:"file_verify_url,omitempty"`

	// 加速域名。
	DomainName *string `json:"domain_name,omitempty"`

	// 校验域名。
	VerifyDomainName *string `json:"verify_domain_name,omitempty"`

	// 文件校验的校验文件名。
	FileVerifyFilename *string `json:"file_verify_filename,omitempty"`

	// 校验值，解析值或者文件内容。
	VerifyContent *string `json:"verify_content,omitempty"`

	// 文件校验域名列表。
	FileVerifyDomains *[]string `json:"file_verify_domains,omitempty"`

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
