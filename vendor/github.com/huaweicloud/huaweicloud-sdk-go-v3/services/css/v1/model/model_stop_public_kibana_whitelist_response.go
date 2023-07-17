package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopPublicKibanaWhitelistResponse Response Object
type StopPublicKibanaWhitelistResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StopPublicKibanaWhitelistResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopPublicKibanaWhitelistResponse struct{}"
	}

	return strings.Join([]string{"StopPublicKibanaWhitelistResponse", string(data)}, " ")
}
