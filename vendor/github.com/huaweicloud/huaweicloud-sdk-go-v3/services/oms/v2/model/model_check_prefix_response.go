package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CheckPrefixResponse Response Object
type CheckPrefixResponse struct {

	// 是否存在
	Exist          *bool `json:"exist,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o CheckPrefixResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckPrefixResponse struct{}"
	}

	return strings.Join([]string{"CheckPrefixResponse", string(data)}, " ")
}
