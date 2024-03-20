package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DownloadRegionCarrierExcelResponse Response Object
type DownloadRegionCarrierExcelResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DownloadRegionCarrierExcelResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DownloadRegionCarrierExcelResponse struct{}"
	}

	return strings.Join([]string{"DownloadRegionCarrierExcelResponse", string(data)}, " ")
}
