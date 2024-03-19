package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateResponseHeaderResponse Response Object
type UpdateResponseHeaderResponse struct {
	Headers *HeaderMap `json:"headers,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateResponseHeaderResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateResponseHeaderResponse struct{}"
	}

	return strings.Join([]string{"UpdateResponseHeaderResponse", string(data)}, " ")
}
