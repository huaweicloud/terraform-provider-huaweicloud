package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneCreateProjectRequest struct {
	Body *KeystoneCreateProjectRequestBody `json:"body,omitempty"`
}

func (o KeystoneCreateProjectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateProjectRequest struct{}"
	}

	return strings.Join([]string{"KeystoneCreateProjectRequest", string(data)}, " ")
}
