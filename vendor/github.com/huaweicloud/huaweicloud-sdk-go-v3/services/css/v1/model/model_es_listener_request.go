package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EsListenerRequest struct {

	// 监听器使用的服务器证书ID。
	DefaultTlsContainerRef string `json:"default_tls_container_ref"`

	// 监听器使用的CA证书ID。如果更新双向认证，则该参数为必选。
	ClientCaTlsContainerRef *string `json:"client_ca_tls_container_ref,omitempty"`
}

func (o EsListenerRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EsListenerRequest struct{}"
	}

	return strings.Join([]string{"EsListenerRequest", string(data)}, " ")
}
