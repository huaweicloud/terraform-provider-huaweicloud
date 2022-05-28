package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ConfirmAssetUploadRequest struct {
	Body *ConfirmAssetUploadReq `json:"body,omitempty"`
}

func (o ConfirmAssetUploadRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfirmAssetUploadRequest struct{}"
	}

	return strings.Join([]string{"ConfirmAssetUploadRequest", string(data)}, " ")
}
