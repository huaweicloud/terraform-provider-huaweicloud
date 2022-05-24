package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateAssetByFileUploadRequest struct {
	Body *CreateAssetByFileUploadReq `json:"body,omitempty"`
}

func (o CreateAssetByFileUploadRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAssetByFileUploadRequest struct{}"
	}

	return strings.Join([]string{"CreateAssetByFileUploadRequest", string(data)}, " ")
}
