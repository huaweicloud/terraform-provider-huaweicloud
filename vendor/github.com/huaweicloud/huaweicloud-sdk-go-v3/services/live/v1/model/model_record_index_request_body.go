package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RecordIndexRequestBody struct {

	// 推流域名
	PublishDomain string `json:"publish_domain"`

	// app名
	App string `json:"app"`

	// 流名
	Stream string `json:"stream"`

	// 开始时间。格式为：YYYY-MM-DDTHH:mm:ssZ（UTC时间），开始时间与结束时间的间隔最大为12小时。
	StartTime *sdktime.SdkTime `json:"start_time"`

	// 结束时间。格式为：YYYY-MM-DDTHH:mm:ssZ（UTC时间），开始时间与结束时间的间隔最大为12小时。结束时间不允许大于当前时间。
	EndTime *sdktime.SdkTime `json:"end_time"`

	// \"m3u8文件在OBS中的储存路径。支持下列字符串的转义   - {publish_domain}   - {app}   - {stream}   - {start_time}   - {end_time} 其中{start_time},{end_time}为返回结果的实际时间。 默认值为Index/{publish_domain}/{app}/{stream}/{stream}-{start_time}-{end_time}\"
	Object *string `json:"object,omitempty"`
}

func (o RecordIndexRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecordIndexRequestBody struct{}"
	}

	return strings.Join([]string{"RecordIndexRequestBody", string(data)}, " ")
}
