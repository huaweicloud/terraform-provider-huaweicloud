package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type RunRecordResponse struct {
	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o RunRecordResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RunRecordResponse struct{}"
	}

	return strings.Join([]string{"RunRecordResponse", string(data)}, " ")
}
