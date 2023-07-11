package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListProvidersResponse Response Object
type ListProvidersResponse struct {

	// 云服务列表
	Providers *[]ProviderResponseBody `json:"providers,omitempty"`

	// 当前支持的云服务总数
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListProvidersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProvidersResponse struct{}"
	}

	return strings.Join([]string{"ListProvidersResponse", string(data)}, " ")
}
