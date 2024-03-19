package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePostgresqlExtensionResponse Response Object
type UpdatePostgresqlExtensionResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdatePostgresqlExtensionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePostgresqlExtensionResponse struct{}"
	}

	return strings.Join([]string{"UpdatePostgresqlExtensionResponse", string(data)}, " ")
}
