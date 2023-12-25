package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCloudTypeResponse Response Object
type ShowCloudTypeResponse struct {
	Body           *[]string `json:"body,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ShowCloudTypeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCloudTypeResponse struct{}"
	}

	return strings.Join([]string{"ShowCloudTypeResponse", string(data)}, " ")
}
