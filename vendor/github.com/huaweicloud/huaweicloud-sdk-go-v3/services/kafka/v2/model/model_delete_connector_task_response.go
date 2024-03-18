package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteConnectorTaskResponse Response Object
type DeleteConnectorTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteConnectorTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteConnectorTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteConnectorTaskResponse", string(data)}, " ")
}
