package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListNotifySmnTopicConfigRequest Request Object
type ListNotifySmnTopicConfigRequest struct {
}

func (o ListNotifySmnTopicConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListNotifySmnTopicConfigRequest struct{}"
	}

	return strings.Join([]string{"ListNotifySmnTopicConfigRequest", string(data)}, " ")
}
