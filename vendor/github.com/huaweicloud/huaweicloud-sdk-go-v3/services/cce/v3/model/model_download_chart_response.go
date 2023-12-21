package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"io"

	"strings"
)

// DownloadChartResponse Response Object
type DownloadChartResponse struct {
	HttpStatusCode int           `json:"-"`
	Body           io.ReadCloser `json:"-" type:"stream"`
}

func (o DownloadChartResponse) Consume(writer io.Writer) (int64, error) {
	written, err := io.Copy(writer, o.Body)
	defer o.Body.Close()

	return written, err
}

func (o DownloadChartResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DownloadChartResponse struct{}"
	}

	return strings.Join([]string{"DownloadChartResponse", string(data)}, " ")
}
