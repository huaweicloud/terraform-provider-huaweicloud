package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UploadMetaDataByUrlRequest struct {
	Body *UploadMetaDataByUrlReq `json:"body,omitempty"`
}

func (o UploadMetaDataByUrlRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UploadMetaDataByUrlRequest struct{}"
	}

	return strings.Join([]string{"UploadMetaDataByUrlRequest", string(data)}, " ")
}
