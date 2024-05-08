package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListOttChannelInfoResponse Response Object
type ListOttChannelInfoResponse struct {

	// 总频道数
	Total *int32 `json:"total,omitempty"`

	// 频道信息
	Channels       *[]CreateOttChannelInfoReq `json:"channels,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o ListOttChannelInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListOttChannelInfoResponse struct{}"
	}

	return strings.Join([]string{"ListOttChannelInfoResponse", string(data)}, " ")
}
