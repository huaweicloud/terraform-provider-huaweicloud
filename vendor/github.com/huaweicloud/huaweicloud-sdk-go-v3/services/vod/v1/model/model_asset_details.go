package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type AssetDetails struct {

	// 媒资ID。
	AssetId string `json:"asset_id"`

	// 媒资状态。 - \"CREATING\"：上传中 - \"FAILED\"：上传失败 - \"CREATED\"：上传成功 - \"PUBLISHED\"：已发布 - \"DELETED\"：已删除
	AssetStatus *AssetDetailsAssetStatus `json:"asset_status,omitempty"`

	// 转码状态。 - \"UN_TRANSCODE\"：未转码 - \"WAITING_TRANSCODE\"：等待转码，排队中 - \"TRANSCODING\"：转码中 - \"TRANSCODE_SUCCEED\"：转码成功 - \"TRANSCODE_FAILED\"：转码失败
	TranscodeStatus *AssetDetailsTranscodeStatus `json:"transcode_status,omitempty"`

	BaseInfo *BaseInfo `json:"base_info,omitempty"`

	TranscodeInfo *TranscodeInfo `json:"transcode_info,omitempty"`
}

func (o AssetDetails) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssetDetails struct{}"
	}

	return strings.Join([]string{"AssetDetails", string(data)}, " ")
}

type AssetDetailsAssetStatus struct {
	value string
}

type AssetDetailsAssetStatusEnum struct {
	CREATING  AssetDetailsAssetStatus
	FAILED    AssetDetailsAssetStatus
	CREATED   AssetDetailsAssetStatus
	PUBLISHED AssetDetailsAssetStatus
	DELETED   AssetDetailsAssetStatus
}

func GetAssetDetailsAssetStatusEnum() AssetDetailsAssetStatusEnum {
	return AssetDetailsAssetStatusEnum{
		CREATING: AssetDetailsAssetStatus{
			value: "CREATING",
		},
		FAILED: AssetDetailsAssetStatus{
			value: "FAILED",
		},
		CREATED: AssetDetailsAssetStatus{
			value: "CREATED",
		},
		PUBLISHED: AssetDetailsAssetStatus{
			value: "PUBLISHED",
		},
		DELETED: AssetDetailsAssetStatus{
			value: "DELETED",
		},
	}
}

func (c AssetDetailsAssetStatus) Value() string {
	return c.value
}

func (c AssetDetailsAssetStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AssetDetailsAssetStatus) UnmarshalJSON(b []byte) error {
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

type AssetDetailsTranscodeStatus struct {
	value string
}

type AssetDetailsTranscodeStatusEnum struct {
	UN_TRANSCODE      AssetDetailsTranscodeStatus
	WAITING_TRANSCODE AssetDetailsTranscodeStatus
	TRANSCODING       AssetDetailsTranscodeStatus
	TRANSCODE_SUCCEED AssetDetailsTranscodeStatus
	TRANSCODE_FAILED  AssetDetailsTranscodeStatus
}

func GetAssetDetailsTranscodeStatusEnum() AssetDetailsTranscodeStatusEnum {
	return AssetDetailsTranscodeStatusEnum{
		UN_TRANSCODE: AssetDetailsTranscodeStatus{
			value: "UN_TRANSCODE",
		},
		WAITING_TRANSCODE: AssetDetailsTranscodeStatus{
			value: "WAITING_TRANSCODE",
		},
		TRANSCODING: AssetDetailsTranscodeStatus{
			value: "TRANSCODING",
		},
		TRANSCODE_SUCCEED: AssetDetailsTranscodeStatus{
			value: "TRANSCODE_SUCCEED",
		},
		TRANSCODE_FAILED: AssetDetailsTranscodeStatus{
			value: "TRANSCODE_FAILED",
		},
	}
}

func (c AssetDetailsTranscodeStatus) Value() string {
	return c.value
}

func (c AssetDetailsTranscodeStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AssetDetailsTranscodeStatus) UnmarshalJSON(b []byte) error {
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
