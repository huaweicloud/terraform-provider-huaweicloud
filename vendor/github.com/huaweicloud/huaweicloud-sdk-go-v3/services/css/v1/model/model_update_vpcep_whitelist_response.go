package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateVpcepWhitelistResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateVpcepWhitelistResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateVpcepWhitelistResponse struct{}"
	}

	return strings.Join([]string{"UpdateVpcepWhitelistResponse", string(data)}, " ")
}
