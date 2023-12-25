package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteReleaseResponse Response Object
type DeleteReleaseResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteReleaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteReleaseResponse struct{}"
	}

	return strings.Join([]string{"DeleteReleaseResponse", string(data)}, " ")
}
