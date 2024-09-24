package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SessionContext 临时安全凭据的属性。
type SessionContext struct {
	Attributes *Attributes `json:"attributes,omitempty"`
}

func (o SessionContext) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SessionContext struct{}"
	}

	return strings.Join([]string{"SessionContext", string(data)}, " ")
}
