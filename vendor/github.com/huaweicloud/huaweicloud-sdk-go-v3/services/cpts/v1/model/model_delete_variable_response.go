package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteVariableResponse Response Object
type DeleteVariableResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteVariableResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteVariableResponse struct{}"
	}

	return strings.Join([]string{"DeleteVariableResponse", string(data)}, " ")
}
