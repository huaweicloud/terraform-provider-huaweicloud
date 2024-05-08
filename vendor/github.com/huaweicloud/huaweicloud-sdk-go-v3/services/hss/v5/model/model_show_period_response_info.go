package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowPeriodResponseInfo struct {

	// 购买时长数值串，多个用逗号分隔，如1,2,3,4,5,6,7,8,9
	PeriodVals *string `json:"period_vals,omitempty"`

	// 购买时长单位   - year ：年   - month ：月   - day ：日
	PeriodUnit *string `json:"period_unit,omitempty"`
}

func (o ShowPeriodResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPeriodResponseInfo struct{}"
	}

	return strings.Join([]string{"ShowPeriodResponseInfo", string(data)}, " ")
}
