package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListHlsConfigResponse Response Object
type ListHlsConfigResponse struct {

	// 推流域名
	PushDomain *string `json:"push_domain,omitempty"`

	// 推流域名APP配置
	Application    *[]PushDomainApplication `json:"application,omitempty"`
	HttpStatusCode int                      `json:"-"`
}

func (o ListHlsConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHlsConfigResponse struct{}"
	}

	return strings.Join([]string{"ListHlsConfigResponse", string(data)}, " ")
}
