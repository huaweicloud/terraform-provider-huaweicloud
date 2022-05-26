package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type RecordContentInfoV2 struct {

	// 直播推流域名
	PublishDomain *string `json:"publish_domain,omitempty"`

	// 录制文件名
	FileName *string `json:"file_name,omitempty"`

	// 应用名
	App *string `json:"app,omitempty"`

	// 录制的流名
	Stream *string `json:"stream,omitempty"`

	// 录制格式flv，hls，mp4
	RecordFormat *RecordContentInfoV2RecordFormat `json:"record_format,omitempty"`

	// 录制类型，CONTINUOUS_RECORD，COMMAND_RECORD，PLAN_RECORD, ON_DEMAND_RECORD。默认CONTINUOUS_RECORD。 - CONTINUOUS_RECORD：持续录制，在该规则类型配置后，只要有流到推送到录制系统，就触发录制。 - COMMAND_RECORD：命令录制，在该规则类型配置后，在流推送到录制系统后，租户需要通过命令控制该流的录制开始和结束。 - PLAN_RECORD：计划录制，在该规则类型配置后，推的流如果在计划录制的时间区间则触发录制。 - ON_DEMAND_RECORD：按需录制，在该规则类型配置后，录制系统收到推流后，需要调用租户提供的接口查询录制规则，并根据规则录制。
	RecordType *RecordContentInfoV2RecordType `json:"record_type,omitempty"`

	ObsAddr *RecordObsFileAddr `json:"obs_addr,omitempty"`

	VodInfo *VodInfoV2 `json:"vod_info,omitempty"`

	// OBS下载地址
	DownloadUrl *string `json:"download_url,omitempty"`

	// 录制开始时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间。对record_type为PLAN_RECORD有效
	StartTime *string `json:"start_time,omitempty"`

	// 录制结束时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间。对record_type为PLAN_RECORD有效
	EndTime *string `json:"end_time,omitempty"`

	// 该录制文件时长，单位为秒
	Duration *int32 `json:"duration,omitempty"`
}

func (o RecordContentInfoV2) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecordContentInfoV2 struct{}"
	}

	return strings.Join([]string{"RecordContentInfoV2", string(data)}, " ")
}

type RecordContentInfoV2RecordFormat struct {
	value string
}

type RecordContentInfoV2RecordFormatEnum struct {
	FLV RecordContentInfoV2RecordFormat
	HLS RecordContentInfoV2RecordFormat
	MP4 RecordContentInfoV2RecordFormat
}

func GetRecordContentInfoV2RecordFormatEnum() RecordContentInfoV2RecordFormatEnum {
	return RecordContentInfoV2RecordFormatEnum{
		FLV: RecordContentInfoV2RecordFormat{
			value: "FLV",
		},
		HLS: RecordContentInfoV2RecordFormat{
			value: "HLS",
		},
		MP4: RecordContentInfoV2RecordFormat{
			value: "MP4",
		},
	}
}

func (c RecordContentInfoV2RecordFormat) Value() string {
	return c.value
}

func (c RecordContentInfoV2RecordFormat) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RecordContentInfoV2RecordFormat) UnmarshalJSON(b []byte) error {
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

type RecordContentInfoV2RecordType struct {
	value string
}

type RecordContentInfoV2RecordTypeEnum struct {
	CONTINUOUS_RECORD RecordContentInfoV2RecordType
	COMMAND_RECORD    RecordContentInfoV2RecordType
	PLAN_RECORD       RecordContentInfoV2RecordType
	ON_DEMAND_RECORD  RecordContentInfoV2RecordType
}

func GetRecordContentInfoV2RecordTypeEnum() RecordContentInfoV2RecordTypeEnum {
	return RecordContentInfoV2RecordTypeEnum{
		CONTINUOUS_RECORD: RecordContentInfoV2RecordType{
			value: "CONTINUOUS_RECORD",
		},
		COMMAND_RECORD: RecordContentInfoV2RecordType{
			value: "COMMAND_RECORD",
		},
		PLAN_RECORD: RecordContentInfoV2RecordType{
			value: "PLAN_RECORD",
		},
		ON_DEMAND_RECORD: RecordContentInfoV2RecordType{
			value: "ON_DEMAND_RECORD",
		},
	}
}

func (c RecordContentInfoV2RecordType) Value() string {
	return c.value
}

func (c RecordContentInfoV2RecordType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RecordContentInfoV2RecordType) UnmarshalJSON(b []byte) error {
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
