package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowApplicationsResponse struct {

	// 资源空间信息列表。
	Applications   *[]ApplicationDto `json:"applications,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o ShowApplicationsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowApplicationsResponse struct{}"
	}

	return strings.Join([]string{"ShowApplicationsResponse", string(data)}, " ")
}
