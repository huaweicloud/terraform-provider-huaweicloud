package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPermissionsResponse Response Object
type ListPermissionsResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListPermissionsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPermissionsResponse struct{}"
	}

	return strings.Join([]string{"ListPermissionsResponse", string(data)}, " ")
}
