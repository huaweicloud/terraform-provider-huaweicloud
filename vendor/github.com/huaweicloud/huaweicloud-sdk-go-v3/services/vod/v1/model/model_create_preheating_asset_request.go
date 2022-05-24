package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreatePreheatingAssetRequest struct {
	Body *CreatePreheatingAssetReq `json:"body,omitempty"`
}

func (o CreatePreheatingAssetRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePreheatingAssetRequest struct{}"
	}

	return strings.Join([]string{"CreatePreheatingAssetRequest", string(data)}, " ")
}
