package model

import (
	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"
	"strings"
)

// Response Object
type CreateRecordIndexResponse struct {

	// 索引文件地址
	IndexUrl *string `json:"index_url,omitempty"`

	// 直播推流域名
	PublishDomain *string `json:"publish_domain,omitempty"`

	// 应用名，如果任意应用填写*。录制规则匹配的时候，优先精确app匹配，如果匹配不到，则匹配*
	App *string `json:"app,omitempty"`

	// 录制的流名，如果任意流名则填写*。录制规则匹配的时候，优先精确stream匹配，如果匹配不到，则匹配*
	Stream *string `json:"stream,omitempty"`

	// 开始时间。格式为：yyyy-MM-ddTHH:mm:ssZ（UTC时间）。(实际视频的开始时间)
	StartTime *sdktime.SdkTime `json:"start_time,omitempty"`

	// 结束时间。格式为：yyyy-MM-ddTHH:mm:ssZ（UTC时间）。(实际视频的结束时间)
	EndTime *sdktime.SdkTime `json:"end_time,omitempty"`

	// 录制时长。单位：秒。
	Duration *int32 `json:"duration,omitempty"`

	// 视频宽。
	Weight *int32 `json:"weight,omitempty"`

	// 视频高。
	Height *int32 `json:"height,omitempty"`

	// OBS Bucket所在RegionID
	Location *CreateRecordIndexResponseLocation `json:"location,omitempty"`

	// 桶名称
	Bucket *string `json:"bucket,omitempty"`

	// m3u8文件路径。默认Index/{publish_domain}/{app}/{stream}-{start_time}-{end_time}
	Object *string `json:"object,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateRecordIndexResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateRecordIndexResponse struct{}"
	}

	return strings.Join([]string{"CreateRecordIndexResponse", string(data)}, " ")
}

type CreateRecordIndexResponseLocation struct {
	value string
}

type CreateRecordIndexResponseLocationEnum struct {
	CN_NORTH_4 CreateRecordIndexResponseLocation
	CN_NORTH_5 CreateRecordIndexResponseLocation
	CN_NORTH_6 CreateRecordIndexResponseLocation
}

func GetCreateRecordIndexResponseLocationEnum() CreateRecordIndexResponseLocationEnum {
	return CreateRecordIndexResponseLocationEnum{
		CN_NORTH_4: CreateRecordIndexResponseLocation{
			value: "cn-north-4",
		},
		CN_NORTH_5: CreateRecordIndexResponseLocation{
			value: "cn-north-5",
		},
		CN_NORTH_6: CreateRecordIndexResponseLocation{
			value: "cn-north-6",
		},
	}
}

func (c CreateRecordIndexResponseLocation) Value() string {
	return c.value
}

func (c CreateRecordIndexResponseLocation) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateRecordIndexResponseLocation) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
