package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UnlockNodeReadonlyStatusResponse Response Object
type UnlockNodeReadonlyStatusResponse struct {

	// 解除结果
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UnlockNodeReadonlyStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnlockNodeReadonlyStatusResponse struct{}"
	}

	return strings.Join([]string{"UnlockNodeReadonlyStatusResponse", string(data)}, " ")
}
