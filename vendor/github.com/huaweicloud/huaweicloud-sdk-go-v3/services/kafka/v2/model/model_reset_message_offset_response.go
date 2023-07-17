package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResetMessageOffsetResponse Response Object
type ResetMessageOffsetResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ResetMessageOffsetResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetMessageOffsetResponse struct{}"
	}

	return strings.Join([]string{"ResetMessageOffsetResponse", string(data)}, " ")
}
