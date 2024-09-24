package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchScanSwrImageRequest Request Object
type BatchScanSwrImageRequest struct {

	// Region ID
	Region *string `json:"region,omitempty"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *BatchScanPrivateImageRequestInfo `json:"body,omitempty"`
}

func (o BatchScanSwrImageRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchScanSwrImageRequest struct{}"
	}

	return strings.Join([]string{"BatchScanSwrImageRequest", string(data)}, " ")
}
