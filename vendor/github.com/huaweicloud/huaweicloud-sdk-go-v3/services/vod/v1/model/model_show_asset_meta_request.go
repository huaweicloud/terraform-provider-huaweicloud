package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ShowAssetMetaRequest struct {

	// 媒资id，最多同时查询10个媒资。
	AssetId *[]string `json:"asset_id,omitempty"`

	// 媒资状态。  取值如下： - UNCREATED：未创建（媒资ID不存在 ） - DELETED：已删除 - CANCELLED：上传取消 - SERVER_ERROR：上传失败（点播服务端故障） - UPLOAD_FAILED：上传失败（向OBS上传失败） - CREATING：创建中 - PUBLISHED：已发布 - TRANSCODING：待发布（转码中） - TRANSCODE_FAILED：待发布（转码失败） - TRANSCODE_SUCCEED：待发布（转码成功） - CREATED：待发布（未转码）
	Status *[]ShowAssetMetaRequestStatus `json:"status,omitempty"`

	// 转码状态  取值如下： - TRANSCODING：转码中 - TRANSCODE_FAILED：转码失败 - TRANSCODE_SUCCEED：转码成功 - UN_TRANSCODE：未转码 - WAITING_TRANSCODE：等待转码
	TranscodeStatus *[]ShowAssetMetaRequestTranscodeStatus `json:"transcodeStatus,omitempty"`

	// 媒资状态。  取值如下： - PUBLISHED：已发布 - CREATED：未发布
	AssetStatus *[]ShowAssetMetaRequestAssetStatus `json:"assetStatus,omitempty"`

	// 起始时间，查询指定“**asset_id**”时，该参数无效。  格式为yyyymmddhhmmss。必须是与时区无关的UTC时间。
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间，查询指定“**asset_id**”时，该参数无效。  格式为yyyymmddhhmmss。必须是与时区无关的UTC时间。
	EndTime *string `json:"end_time,omitempty"`

	// 分类ID。
	CategoryId *int32 `json:"category_id,omitempty"`

	// 媒资标签。  单个标签不超过16个字节，最多不超过16个标签。  多个用逗号分隔，UTF8编码。
	Tags *string `json:"tags,omitempty"`

	// 在媒资标题、描述中模糊查询的字符串。
	QueryString *string `json:"query_string,omitempty"`

	// 分页编号，查询指定“asset_id”时，该参数无效。  默认值：0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数，查询指定“**asset_id**”时，该参数无效。  取值范围：[1,100]。  默认值：10。
	Size *int32 `json:"size,omitempty"`
}

func (o ShowAssetMetaRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAssetMetaRequest struct{}"
	}

	return strings.Join([]string{"ShowAssetMetaRequest", string(data)}, " ")
}

type ShowAssetMetaRequestStatus struct {
	value string
}

type ShowAssetMetaRequestStatusEnum struct {
	UNCREATED         ShowAssetMetaRequestStatus
	DELETED           ShowAssetMetaRequestStatus
	CANCELLED         ShowAssetMetaRequestStatus
	SERVER_ERROR      ShowAssetMetaRequestStatus
	UPLOAD_FAILED     ShowAssetMetaRequestStatus
	CREATING          ShowAssetMetaRequestStatus
	PUBLISHED         ShowAssetMetaRequestStatus
	WAITING_TRANSCODE ShowAssetMetaRequestStatus
	TRANSCODING       ShowAssetMetaRequestStatus
	TRANSCODE_FAILED  ShowAssetMetaRequestStatus
	TRANSCODE_SUCCEED ShowAssetMetaRequestStatus
	CREATED           ShowAssetMetaRequestStatus
}

func GetShowAssetMetaRequestStatusEnum() ShowAssetMetaRequestStatusEnum {
	return ShowAssetMetaRequestStatusEnum{
		UNCREATED: ShowAssetMetaRequestStatus{
			value: "UNCREATED",
		},
		DELETED: ShowAssetMetaRequestStatus{
			value: "DELETED",
		},
		CANCELLED: ShowAssetMetaRequestStatus{
			value: "CANCELLED",
		},
		SERVER_ERROR: ShowAssetMetaRequestStatus{
			value: "SERVER_ERROR",
		},
		UPLOAD_FAILED: ShowAssetMetaRequestStatus{
			value: "UPLOAD_FAILED",
		},
		CREATING: ShowAssetMetaRequestStatus{
			value: "CREATING",
		},
		PUBLISHED: ShowAssetMetaRequestStatus{
			value: "PUBLISHED",
		},
		WAITING_TRANSCODE: ShowAssetMetaRequestStatus{
			value: "WAITING_TRANSCODE",
		},
		TRANSCODING: ShowAssetMetaRequestStatus{
			value: "TRANSCODING",
		},
		TRANSCODE_FAILED: ShowAssetMetaRequestStatus{
			value: "TRANSCODE_FAILED",
		},
		TRANSCODE_SUCCEED: ShowAssetMetaRequestStatus{
			value: "TRANSCODE_SUCCEED",
		},
		CREATED: ShowAssetMetaRequestStatus{
			value: "CREATED",
		},
	}
}

func (c ShowAssetMetaRequestStatus) Value() string {
	return c.value
}

func (c ShowAssetMetaRequestStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowAssetMetaRequestStatus) UnmarshalJSON(b []byte) error {
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

type ShowAssetMetaRequestTranscodeStatus struct {
	value string
}

type ShowAssetMetaRequestTranscodeStatusEnum struct {
	TRANSCODING       ShowAssetMetaRequestTranscodeStatus
	TRANSCODE_FAILED  ShowAssetMetaRequestTranscodeStatus
	TRANSCODE_SUCCEED ShowAssetMetaRequestTranscodeStatus
	UN_TRANSCODE      ShowAssetMetaRequestTranscodeStatus
	WAITING_TRANSCODE ShowAssetMetaRequestTranscodeStatus
}

func GetShowAssetMetaRequestTranscodeStatusEnum() ShowAssetMetaRequestTranscodeStatusEnum {
	return ShowAssetMetaRequestTranscodeStatusEnum{
		TRANSCODING: ShowAssetMetaRequestTranscodeStatus{
			value: "TRANSCODING",
		},
		TRANSCODE_FAILED: ShowAssetMetaRequestTranscodeStatus{
			value: "TRANSCODE_FAILED",
		},
		TRANSCODE_SUCCEED: ShowAssetMetaRequestTranscodeStatus{
			value: "TRANSCODE_SUCCEED",
		},
		UN_TRANSCODE: ShowAssetMetaRequestTranscodeStatus{
			value: "UN_TRANSCODE",
		},
		WAITING_TRANSCODE: ShowAssetMetaRequestTranscodeStatus{
			value: "WAITING_TRANSCODE",
		},
	}
}

func (c ShowAssetMetaRequestTranscodeStatus) Value() string {
	return c.value
}

func (c ShowAssetMetaRequestTranscodeStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowAssetMetaRequestTranscodeStatus) UnmarshalJSON(b []byte) error {
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

type ShowAssetMetaRequestAssetStatus struct {
	value string
}

type ShowAssetMetaRequestAssetStatusEnum struct {
	PUBLISHED ShowAssetMetaRequestAssetStatus
	CREATED   ShowAssetMetaRequestAssetStatus
}

func GetShowAssetMetaRequestAssetStatusEnum() ShowAssetMetaRequestAssetStatusEnum {
	return ShowAssetMetaRequestAssetStatusEnum{
		PUBLISHED: ShowAssetMetaRequestAssetStatus{
			value: "PUBLISHED",
		},
		CREATED: ShowAssetMetaRequestAssetStatus{
			value: "CREATED",
		},
	}
}

func (c ShowAssetMetaRequestAssetStatus) Value() string {
	return c.value
}

func (c ShowAssetMetaRequestAssetStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowAssetMetaRequestAssetStatus) UnmarshalJSON(b []byte) error {
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
