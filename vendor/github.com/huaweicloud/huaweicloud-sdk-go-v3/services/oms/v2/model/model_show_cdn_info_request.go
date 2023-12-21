package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCdnInfoRequest Request Object
type ShowCdnInfoRequest struct {
	Body *ShowCdnInfoReq `json:"body,omitempty"`
}

func (o ShowCdnInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCdnInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowCdnInfoRequest", string(data)}, " ")
}
