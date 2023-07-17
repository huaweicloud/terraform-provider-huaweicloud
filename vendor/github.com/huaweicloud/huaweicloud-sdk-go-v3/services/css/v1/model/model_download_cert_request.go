package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DownloadCertRequest Request Object
type DownloadCertRequest struct {
}

func (o DownloadCertRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DownloadCertRequest struct{}"
	}

	return strings.Join([]string{"DownloadCertRequest", string(data)}, " ")
}
