package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreatePartitionReq struct {

	// 期望调整分区后的数量，必须大于当前分区数量，小于等于 [100](tag:hc,hk,hws,hws_hk,otc,hws_ocb,ctc,sbc,hk_sbc,g42,tm)[20](tag:cmcc)。
	Partition *int32 `json:"partition,omitempty"`
}

func (o CreatePartitionReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePartitionReq struct{}"
	}

	return strings.Join([]string{"CreatePartitionReq", string(data)}, " ")
}
