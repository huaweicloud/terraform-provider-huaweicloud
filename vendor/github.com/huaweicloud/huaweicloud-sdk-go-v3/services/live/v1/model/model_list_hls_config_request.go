package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListHlsConfigRequest Request Object
type ListHlsConfigRequest struct {

	// 推流域名
	PushDomain string `json:"push_domain"`
}

func (o ListHlsConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHlsConfigRequest struct{}"
	}

	return strings.Join([]string{"ListHlsConfigRequest", string(data)}, " ")
}
