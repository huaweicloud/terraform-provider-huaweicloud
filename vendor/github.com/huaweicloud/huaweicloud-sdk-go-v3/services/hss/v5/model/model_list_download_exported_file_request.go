package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDownloadExportedFileRequest Request Object
type ListDownloadExportedFileRequest struct {

	// 文件id
	FileId string `json:"file_id"`

	// Region Id
	Region string `json:"region"`

	// 企业项目ID
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ListDownloadExportedFileRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDownloadExportedFileRequest struct{}"
	}

	return strings.Join([]string{"ListDownloadExportedFileRequest", string(data)}, " ")
}
