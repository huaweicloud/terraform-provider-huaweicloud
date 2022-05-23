package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateAccessCodeResponse struct {

	// 接入名，随机生成8位字符串
	AccessKey *string `json:"access_key,omitempty"`

	// 接入凭证。
	AccessCode     *string `json:"access_code,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateAccessCodeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAccessCodeResponse struct{}"
	}

	return strings.Join([]string{"CreateAccessCodeResponse", string(data)}, " ")
}
