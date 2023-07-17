package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"io"

	"strings"
)

// DownloadCertResponse Response Object
type DownloadCertResponse struct {
	HttpStatusCode int           `json:"-"`
	Body           io.ReadCloser `json:"-" type:"stream"`
}

func (o DownloadCertResponse) Consume(writer io.Writer) (int64, error) {
	written, err := io.Copy(writer, o.Body)
	defer o.Body.Close()

	return written, err
}

func (o DownloadCertResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DownloadCertResponse struct{}"
	}

	return strings.Join([]string{"DownloadCertResponse", string(data)}, " ")
}
