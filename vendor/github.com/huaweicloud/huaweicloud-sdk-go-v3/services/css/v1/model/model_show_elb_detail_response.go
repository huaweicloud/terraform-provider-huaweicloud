package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowElbDetailResponse Response Object
type ShowElbDetailResponse struct {

	// 服务器证书名称。
	ServerCertName *string `json:"serverCertName,omitempty"`

	// 服务器证书ID。
	ServerCertId *string `json:"serverCertId,omitempty"`

	// ca证书名称。
	CacertName *string `json:"cacertName,omitempty"`

	// ca证书ID。
	CacertId *string `json:"cacertId,omitempty"`

	// elb开关信息。
	ElbEnable *bool `json:"elb_enable,omitempty"`

	// 认证方式。
	AuthenticationType *string `json:"authentication_type,omitempty"`

	LoadBalancer *EsLoadBalancerResource `json:"loadBalancer,omitempty"`

	Healthmonitors *EsHealthmonitorsResource `json:"healthmonitors,omitempty"`
	HttpStatusCode int                       `json:"-"`
}

func (o ShowElbDetailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowElbDetailResponse struct{}"
	}

	return strings.Join([]string{"ShowElbDetailResponse", string(data)}, " ")
}
