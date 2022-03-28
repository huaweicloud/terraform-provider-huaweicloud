package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteAgencyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteAgencyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAgencyResponse struct{}"
	}

	return strings.Join([]string{"DeleteAgencyResponse", string(data)}, " ")
}
