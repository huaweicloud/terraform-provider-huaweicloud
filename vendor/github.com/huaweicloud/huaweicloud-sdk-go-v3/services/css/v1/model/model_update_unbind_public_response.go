package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateUnbindPublicResponse struct {

	// 操作行为。固定为：unbindZone，表示解绑成功。
	Action         *string `json:"action,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateUnbindPublicResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateUnbindPublicResponse struct{}"
	}

	return strings.Join([]string{"UpdateUnbindPublicResponse", string(data)}, " ")
}
