package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowAgencyRequest struct {

	// 待查询的委托ID，获取方式请参见：[获取委托名、委托ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	AgencyId string `json:"agency_id"`
}

func (o ShowAgencyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAgencyRequest struct{}"
	}

	return strings.Join([]string{"ShowAgencyRequest", string(data)}, " ")
}
