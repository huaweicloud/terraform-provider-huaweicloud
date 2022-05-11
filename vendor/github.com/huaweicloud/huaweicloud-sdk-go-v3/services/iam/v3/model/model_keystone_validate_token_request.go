package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneValidateTokenRequest struct {

	// 待校验的token。
	XSubjectToken string `json:"X-Subject-Token"`

	// 如果设置该参数，返回的响应体中将不显示catalog信息。任何非空字符串都将解释为true，并使该字段生效。
	Nocatalog *string `json:"nocatalog,omitempty"`
}

func (o KeystoneValidateTokenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneValidateTokenRequest struct{}"
	}

	return strings.Join([]string{"KeystoneValidateTokenRequest", string(data)}, " ")
}
