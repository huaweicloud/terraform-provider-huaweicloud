package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePublicBandWidthRequest Request Object
type UpdatePublicBandWidthRequest struct {

	// 指定修改公网访问带宽集群ID。
	ClusterId string `json:"cluster_id"`

	Body *BindPublicReqEipReq `json:"body,omitempty"`
}

func (o UpdatePublicBandWidthRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePublicBandWidthRequest struct{}"
	}

	return strings.Join([]string{"UpdatePublicBandWidthRequest", string(data)}, " ")
}
