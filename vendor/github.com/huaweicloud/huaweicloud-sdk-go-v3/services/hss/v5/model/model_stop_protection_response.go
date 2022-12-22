package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type StopProtectionResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StopProtectionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopProtectionResponse struct{}"
	}

	return strings.Join([]string{"StopProtectionResponse", string(data)}, " ")
}
