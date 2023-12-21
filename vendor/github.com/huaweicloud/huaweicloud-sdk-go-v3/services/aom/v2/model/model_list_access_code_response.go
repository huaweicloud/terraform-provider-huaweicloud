package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAccessCodeResponse Response Object
type ListAccessCodeResponse struct {

	// accessCodes。
	AccessCodes    *[]AccessCodeModel `json:"access_codes,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ListAccessCodeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAccessCodeResponse struct{}"
	}

	return strings.Join([]string{"ListAccessCodeResponse", string(data)}, " ")
}
