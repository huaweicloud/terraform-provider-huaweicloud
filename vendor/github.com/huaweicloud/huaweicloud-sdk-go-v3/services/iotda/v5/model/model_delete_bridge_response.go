package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteBridgeResponse Response Object
type DeleteBridgeResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteBridgeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteBridgeResponse struct{}"
	}

	return strings.Join([]string{"DeleteBridgeResponse", string(data)}, " ")
}
