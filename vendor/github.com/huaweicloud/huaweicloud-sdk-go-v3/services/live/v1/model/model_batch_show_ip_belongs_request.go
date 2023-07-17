package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchShowIpBelongsRequest Request Object
type BatchShowIpBelongsRequest struct {

	// IP地址列表，最多20个。
	Ip []string `json:"ip"`
}

func (o BatchShowIpBelongsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchShowIpBelongsRequest struct{}"
	}

	return strings.Join([]string{"BatchShowIpBelongsRequest", string(data)}, " ")
}
