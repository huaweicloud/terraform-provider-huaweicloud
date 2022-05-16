package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// scope详细信息
type GetIdTokenScopeDomainOrProjectBody struct {

	// domain id 或者 project id，与name字段至少存在一个。
	Id *string `json:"id,omitempty"`

	// domain name 或者 project name，与id字段至少存在一个。
	Name *string `json:"name,omitempty"`
}

func (o GetIdTokenScopeDomainOrProjectBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetIdTokenScopeDomainOrProjectBody struct{}"
	}

	return strings.Join([]string{"GetIdTokenScopeDomainOrProjectBody", string(data)}, " ")
}
