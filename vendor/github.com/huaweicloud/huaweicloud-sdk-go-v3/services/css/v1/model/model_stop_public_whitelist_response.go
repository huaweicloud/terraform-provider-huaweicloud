package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type StopPublicWhitelistResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StopPublicWhitelistResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopPublicWhitelistResponse struct{}"
	}

	return strings.Join([]string{"StopPublicWhitelistResponse", string(data)}, " ")
}
