package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateUserInformationResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateUserInformationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateUserInformationResponse struct{}"
	}

	return strings.Join([]string{"UpdateUserInformationResponse", string(data)}, " ")
}
