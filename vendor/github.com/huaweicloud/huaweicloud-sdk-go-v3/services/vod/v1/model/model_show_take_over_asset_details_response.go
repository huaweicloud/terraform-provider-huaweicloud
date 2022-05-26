package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Response Object
type ShowTakeOverAssetDetailsResponse struct {

	// 媒资ID。
	AssetId *string `json:"asset_id,omitempty"`

	// 媒资状态。 - \"CREATING\"：上传中 - \"FAILED\"：上传失败 - \"CREATED\"：上传成功 - \"PUBLISHED\"：已发布 - \"DELETED\"：已删除
	AssetStatus *ShowTakeOverAssetDetailsResponseAssetStatus `json:"asset_status,omitempty"`

	// 转码状态。 - \"UN_TRANSCODE\"：未转码 - \"WAITING_TRANSCODE\"：等待转码，排队中 - \"TRANSCODING\"：转码中 - \"TRANSCODE_SUCCEED\"：转码成功 - \"TRANSCODE_FAILED\"：转码失败
	TranscodeStatus *ShowTakeOverAssetDetailsResponseTranscodeStatus `json:"transcode_status,omitempty"`

	BaseInfo *BaseInfo `json:"base_info,omitempty"`

	TranscodeInfo  *TranscodeInfo `json:"transcode_info,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ShowTakeOverAssetDetailsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTakeOverAssetDetailsResponse struct{}"
	}

	return strings.Join([]string{"ShowTakeOverAssetDetailsResponse", string(data)}, " ")
}

type ShowTakeOverAssetDetailsResponseAssetStatus struct {
	value string
}

type ShowTakeOverAssetDetailsResponseAssetStatusEnum struct {
	CREATING  ShowTakeOverAssetDetailsResponseAssetStatus
	FAILED    ShowTakeOverAssetDetailsResponseAssetStatus
	CREATED   ShowTakeOverAssetDetailsResponseAssetStatus
	PUBLISHED ShowTakeOverAssetDetailsResponseAssetStatus
	DELETED   ShowTakeOverAssetDetailsResponseAssetStatus
}

func GetShowTakeOverAssetDetailsResponseAssetStatusEnum() ShowTakeOverAssetDetailsResponseAssetStatusEnum {
	return ShowTakeOverAssetDetailsResponseAssetStatusEnum{
		CREATING: ShowTakeOverAssetDetailsResponseAssetStatus{
			value: "CREATING",
		},
		FAILED: ShowTakeOverAssetDetailsResponseAssetStatus{
			value: "FAILED",
		},
		CREATED: ShowTakeOverAssetDetailsResponseAssetStatus{
			value: "CREATED",
		},
		PUBLISHED: ShowTakeOverAssetDetailsResponseAssetStatus{
			value: "PUBLISHED",
		},
		DELETED: ShowTakeOverAssetDetailsResponseAssetStatus{
			value: "DELETED",
		},
	}
}

func (c ShowTakeOverAssetDetailsResponseAssetStatus) Value() string {
	return c.value
}

func (c ShowTakeOverAssetDetailsResponseAssetStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTakeOverAssetDetailsResponseAssetStatus) UnmarshalJSON(b []byte) error {
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

type ShowTakeOverAssetDetailsResponseTranscodeStatus struct {
	value string
}

type ShowTakeOverAssetDetailsResponseTranscodeStatusEnum struct {
	UN_TRANSCODE      ShowTakeOverAssetDetailsResponseTranscodeStatus
	WAITING_TRANSCODE ShowTakeOverAssetDetailsResponseTranscodeStatus
	TRANSCODING       ShowTakeOverAssetDetailsResponseTranscodeStatus
	TRANSCODE_SUCCEED ShowTakeOverAssetDetailsResponseTranscodeStatus
	TRANSCODE_FAILED  ShowTakeOverAssetDetailsResponseTranscodeStatus
}

func GetShowTakeOverAssetDetailsResponseTranscodeStatusEnum() ShowTakeOverAssetDetailsResponseTranscodeStatusEnum {
	return ShowTakeOverAssetDetailsResponseTranscodeStatusEnum{
		UN_TRANSCODE: ShowTakeOverAssetDetailsResponseTranscodeStatus{
			value: "UN_TRANSCODE",
		},
		WAITING_TRANSCODE: ShowTakeOverAssetDetailsResponseTranscodeStatus{
			value: "WAITING_TRANSCODE",
		},
		TRANSCODING: ShowTakeOverAssetDetailsResponseTranscodeStatus{
			value: "TRANSCODING",
		},
		TRANSCODE_SUCCEED: ShowTakeOverAssetDetailsResponseTranscodeStatus{
			value: "TRANSCODE_SUCCEED",
		},
		TRANSCODE_FAILED: ShowTakeOverAssetDetailsResponseTranscodeStatus{
			value: "TRANSCODE_FAILED",
		},
	}
}

func (c ShowTakeOverAssetDetailsResponseTranscodeStatus) Value() string {
	return c.value
}

func (c ShowTakeOverAssetDetailsResponseTranscodeStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTakeOverAssetDetailsResponseTranscodeStatus) UnmarshalJSON(b []byte) error {
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
