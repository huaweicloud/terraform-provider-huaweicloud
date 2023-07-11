package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClusterPublicEip 公网带宽信息。
type CreateClusterPublicEip struct {
	BandWidth *CreateClusterPublicEipSize `json:"bandWidth"`
}

func (o CreateClusterPublicEip) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterPublicEip struct{}"
	}

	return strings.Join([]string{"CreateClusterPublicEip", string(data)}, " ")
}
