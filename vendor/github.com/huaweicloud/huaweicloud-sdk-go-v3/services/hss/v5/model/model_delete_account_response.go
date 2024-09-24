package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteAccountResponse Response Object
type DeleteAccountResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteAccountResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAccountResponse struct{}"
	}

	return strings.Join([]string{"DeleteAccountResponse", string(data)}, " ")
}
