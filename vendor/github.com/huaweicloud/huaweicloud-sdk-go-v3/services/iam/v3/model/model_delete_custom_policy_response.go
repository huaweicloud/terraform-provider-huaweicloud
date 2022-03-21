package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteCustomPolicyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteCustomPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteCustomPolicyResponse struct{}"
	}

	return strings.Join([]string{"DeleteCustomPolicyResponse", string(data)}, " ")
}
