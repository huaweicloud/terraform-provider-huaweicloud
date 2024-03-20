package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCertificatesHttpsInfoResponse Response Object
type ShowCertificatesHttpsInfoResponse struct {

	// 查询结果总数
	Total *int32 `json:"total,omitempty"`

	// https配置。
	Https *[]HttpsDetail `json:"https,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowCertificatesHttpsInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCertificatesHttpsInfoResponse struct{}"
	}

	return strings.Join([]string{"ShowCertificatesHttpsInfoResponse", string(data)}, " ")
}
