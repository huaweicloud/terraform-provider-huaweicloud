package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"io"

	"strings"
)

// DownloadAutopilotChartResponse Response Object
type DownloadAutopilotChartResponse struct {
	HttpStatusCode int           `json:"-"`
	Body           io.ReadCloser `json:"-" type:"stream"`
}

func (o DownloadAutopilotChartResponse) Consume(writer io.Writer) (int64, error) {
	written, err := io.Copy(writer, o.Body)
	defer o.Body.Close()

	return written, err
}

func (o DownloadAutopilotChartResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DownloadAutopilotChartResponse struct{}"
	}

	return strings.Join([]string{"DownloadAutopilotChartResponse", string(data)}, " ")
}
