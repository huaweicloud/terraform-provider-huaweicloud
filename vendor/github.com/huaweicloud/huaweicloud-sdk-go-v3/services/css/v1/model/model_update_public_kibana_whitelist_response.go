package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdatePublicKibanaWhitelistResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdatePublicKibanaWhitelistResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePublicKibanaWhitelistResponse struct{}"
	}

	return strings.Join([]string{"UpdatePublicKibanaWhitelistResponse", string(data)}, " ")
}
