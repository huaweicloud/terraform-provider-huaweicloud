package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowDomainRequest struct {

	// 直播域名，如果不设置此字段，则返回租户所有的域名信息
	Domain *string `json:"domain,omitempty"`
}

func (o ShowDomainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainRequest struct{}"
	}

	return strings.Join([]string{"ShowDomainRequest", string(data)}, " ")
}
