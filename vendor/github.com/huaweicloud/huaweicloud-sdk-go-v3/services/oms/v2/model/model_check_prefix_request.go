package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CheckPrefixRequest Request Object
type CheckPrefixRequest struct {
	Body *CheckPrefixReq `json:"body,omitempty"`
}

func (o CheckPrefixRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckPrefixRequest struct{}"
	}

	return strings.Join([]string{"CheckPrefixRequest", string(data)}, " ")
}
