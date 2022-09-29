package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteInstanceResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteInstanceResponse struct{}"
	}

	return strings.Join([]string{"DeleteInstanceResponse", string(data)}, " ")
}
