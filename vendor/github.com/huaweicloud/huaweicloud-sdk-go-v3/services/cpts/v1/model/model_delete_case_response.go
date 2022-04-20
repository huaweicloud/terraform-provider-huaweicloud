package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteCaseResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteCaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteCaseResponse struct{}"
	}

	return strings.Join([]string{"DeleteCaseResponse", string(data)}, " ")
}
