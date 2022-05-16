package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteDomainRequest struct {

	// 直播域名
	Domain string `json:"domain"`
}

func (o DeleteDomainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDomainRequest struct{}"
	}

	return strings.Join([]string{"DeleteDomainRequest", string(data)}, " ")
}
