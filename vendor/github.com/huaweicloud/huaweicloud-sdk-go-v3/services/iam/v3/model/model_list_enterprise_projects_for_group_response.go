package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListEnterpriseProjectsForGroupResponse struct {

	// 企业项目信息。
	EnterpriseProjects *[]ListEnterpriseProjectsResDetail `json:"enterprise-projects,omitempty"`
	HttpStatusCode     int                                `json:"-"`
}

func (o ListEnterpriseProjectsForGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEnterpriseProjectsForGroupResponse struct{}"
	}

	return strings.Join([]string{"ListEnterpriseProjectsForGroupResponse", string(data)}, " ")
}
