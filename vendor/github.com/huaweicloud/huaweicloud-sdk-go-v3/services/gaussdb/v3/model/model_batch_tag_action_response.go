package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type BatchTagActionResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchTagActionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchTagActionResponse struct{}"
	}

	return strings.Join([]string{"BatchTagActionResponse", string(data)}, " ")
}
