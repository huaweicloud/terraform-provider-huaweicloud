package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateCopyStateResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateCopyStateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCopyStateResponse struct{}"
	}

	return strings.Join([]string{"UpdateCopyStateResponse", string(data)}, " ")
}
