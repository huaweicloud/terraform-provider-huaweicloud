package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartPublicWhitelistResponse Response Object
type StartPublicWhitelistResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StartPublicWhitelistResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartPublicWhitelistResponse struct{}"
	}

	return strings.Join([]string{"StartPublicWhitelistResponse", string(data)}, " ")
}
