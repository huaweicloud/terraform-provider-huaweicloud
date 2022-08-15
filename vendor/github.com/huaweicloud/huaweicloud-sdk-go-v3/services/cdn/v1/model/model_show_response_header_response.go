package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowResponseHeaderResponse struct {
	Headers        *HeaderMap `json:"headers,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o ShowResponseHeaderResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowResponseHeaderResponse struct{}"
	}

	return strings.Join([]string{"ShowResponseHeaderResponse", string(data)}, " ")
}
