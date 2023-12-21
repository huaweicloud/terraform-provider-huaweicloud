package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RunImageSynchronizeResponse Response Object
type RunImageSynchronizeResponse struct {

	// 错误编码
	ErrorCode *int32 `json:"error_code,omitempty"`

	// 错误描述
	ErrorDescription *string `json:"error_description,omitempty"`
	HttpStatusCode   int     `json:"-"`
}

func (o RunImageSynchronizeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RunImageSynchronizeResponse struct{}"
	}

	return strings.Join([]string{"RunImageSynchronizeResponse", string(data)}, " ")
}
