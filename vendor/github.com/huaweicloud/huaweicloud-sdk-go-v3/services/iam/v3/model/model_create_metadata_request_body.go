package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type CreateMetadataRequestBody struct {

	// 用户所属账号ID。
	DomainId string `json:"domain_id"`

	// 该字段为标识租户来源字段，默认为空。
	XaccountType string `json:"xaccount_type"`

	// 该字段为用户IdP服务器的Metadata文件的内容。
	Metadata string `json:"metadata"`
}

func (o CreateMetadataRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMetadataRequestBody struct{}"
	}

	return strings.Join([]string{"CreateMetadataRequestBody", string(data)}, " ")
}
