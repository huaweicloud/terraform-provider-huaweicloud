package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchShowIpBelongsResponse Response Object
type BatchShowIpBelongsResponse struct {

	// IP归属信息列表。
	CdnIps         *[]CdnIp `json:"cdn_ips,omitempty"`
	HttpStatusCode int      `json:"-"`
}

func (o BatchShowIpBelongsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchShowIpBelongsResponse struct{}"
	}

	return strings.Join([]string{"BatchShowIpBelongsResponse", string(data)}, " ")
}
