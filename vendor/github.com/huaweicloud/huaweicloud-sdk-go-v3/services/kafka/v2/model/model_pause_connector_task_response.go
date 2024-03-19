package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PauseConnectorTaskResponse Response Object
type PauseConnectorTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o PauseConnectorTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PauseConnectorTaskResponse struct{}"
	}

	return strings.Join([]string{"PauseConnectorTaskResponse", string(data)}, " ")
}
