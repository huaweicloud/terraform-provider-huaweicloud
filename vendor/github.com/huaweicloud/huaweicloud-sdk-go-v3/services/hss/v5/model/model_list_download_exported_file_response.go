package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDownloadExportedFileResponse Response Object
type ListDownloadExportedFileResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListDownloadExportedFileResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDownloadExportedFileResponse struct{}"
	}

	return strings.Join([]string{"ListDownloadExportedFileResponse", string(data)}, " ")
}
