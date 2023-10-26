package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ClearPrivateKeyResponse Response Object
type ClearPrivateKeyResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ClearPrivateKeyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClearPrivateKeyResponse struct{}"
	}

	return strings.Join([]string{"ClearPrivateKeyResponse", string(data)}, " ")
}
