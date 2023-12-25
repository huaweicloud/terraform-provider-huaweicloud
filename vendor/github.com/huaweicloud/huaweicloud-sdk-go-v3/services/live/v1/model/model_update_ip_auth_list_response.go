package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateIpAuthListResponse Response Object
type UpdateIpAuthListResponse struct {
	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateIpAuthListResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateIpAuthListResponse struct{}"
	}

	return strings.Join([]string{"UpdateIpAuthListResponse", string(data)}, " ")
}
