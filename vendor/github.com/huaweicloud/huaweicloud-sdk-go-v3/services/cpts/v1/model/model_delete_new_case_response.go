package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteNewCaseResponse Response Object
type DeleteNewCaseResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteNewCaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteNewCaseResponse struct{}"
	}

	return strings.Join([]string{"DeleteNewCaseResponse", string(data)}, " ")
}
