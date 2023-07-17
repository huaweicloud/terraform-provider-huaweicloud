package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateVpcepConnectionResponse Response Object
type UpdateVpcepConnectionResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateVpcepConnectionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateVpcepConnectionResponse struct{}"
	}

	return strings.Join([]string{"UpdateVpcepConnectionResponse", string(data)}, " ")
}
