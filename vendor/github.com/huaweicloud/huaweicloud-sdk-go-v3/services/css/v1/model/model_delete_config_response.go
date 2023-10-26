package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteConfigResponse Response Object
type DeleteConfigResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteConfigResponse struct{}"
	}

	return strings.Join([]string{"DeleteConfigResponse", string(data)}, " ")
}
