package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClusterResponse Response Object
type ShowClusterResponse struct {
	Cluster        *ShowClusterRespCluster `json:"cluster,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o ShowClusterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterResponse struct{}"
	}

	return strings.Join([]string{"ShowClusterResponse", string(data)}, " ")
}
