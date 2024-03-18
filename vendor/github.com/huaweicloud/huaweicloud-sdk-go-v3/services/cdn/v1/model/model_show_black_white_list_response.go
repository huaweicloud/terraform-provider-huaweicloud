package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBlackWhiteListResponse Response Object
type ShowBlackWhiteListResponse struct {

	// IP黑白名单类型（0：关闭IP黑白名单功能，1：黑名单，2：白名单）。
	Type *int32 `json:"type,omitempty"`

	// IP黑白名单列表。
	IpList *[]string `json:"ip_list,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowBlackWhiteListResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBlackWhiteListResponse struct{}"
	}

	return strings.Join([]string{"ShowBlackWhiteListResponse", string(data)}, " ")
}
