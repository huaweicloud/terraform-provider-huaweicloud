package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateBandwidthPolicyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateBandwidthPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBandwidthPolicyResponse struct{}"
	}

	return strings.Join([]string{"UpdateBandwidthPolicyResponse", string(data)}, " ")
}
