package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DownloadStatisticsExcelResponse Response Object
type DownloadStatisticsExcelResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DownloadStatisticsExcelResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DownloadStatisticsExcelResponse struct{}"
	}

	return strings.Join([]string{"DownloadStatisticsExcelResponse", string(data)}, " ")
}
