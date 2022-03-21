package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneListRegionsRequest struct {
}

func (o KeystoneListRegionsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListRegionsRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListRegionsRequest", string(data)}, " ")
}
