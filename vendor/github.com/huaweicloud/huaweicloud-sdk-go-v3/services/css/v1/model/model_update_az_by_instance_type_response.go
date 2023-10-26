package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAzByInstanceTypeResponse Response Object
type UpdateAzByInstanceTypeResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateAzByInstanceTypeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAzByInstanceTypeResponse struct{}"
	}

	return strings.Join([]string{"UpdateAzByInstanceTypeResponse", string(data)}, " ")
}
