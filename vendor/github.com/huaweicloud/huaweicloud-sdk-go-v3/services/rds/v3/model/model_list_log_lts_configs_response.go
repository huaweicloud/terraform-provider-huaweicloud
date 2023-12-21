package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListLogLtsConfigsResponse Response Object
type ListLogLtsConfigsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ListLogLtsConfigsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListLogLtsConfigsResponse struct{}"
	}

	return strings.Join([]string{"ListLogLtsConfigsResponse", string(data)}, " ")
}
