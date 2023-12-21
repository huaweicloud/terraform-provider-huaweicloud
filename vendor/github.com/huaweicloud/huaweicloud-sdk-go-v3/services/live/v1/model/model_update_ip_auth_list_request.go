package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateIpAuthListRequest Request Object
type UpdateIpAuthListRequest struct {
	Body *IpAuthInfo `json:"body,omitempty"`
}

func (o UpdateIpAuthListRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateIpAuthListRequest struct{}"
	}

	return strings.Join([]string{"UpdateIpAuthListRequest", string(data)}, " ")
}
