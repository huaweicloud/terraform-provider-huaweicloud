package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateReferResponse Response Object
type UpdateReferResponse struct {
	Referer *RefererRsp `json:"referer,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateReferResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateReferResponse struct{}"
	}

	return strings.Join([]string{"UpdateReferResponse", string(data)}, " ")
}
