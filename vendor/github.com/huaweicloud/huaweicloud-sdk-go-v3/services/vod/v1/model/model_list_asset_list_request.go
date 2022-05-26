package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListAssetListRequest struct {

	// 媒资ID，最多同时查询10个。
	AssetId *[]string `json:"asset_id,omitempty"`

	// 媒资状态，同时查询多个状态的媒资。  取值如下： - CREATING：上传中 - FAILED：上传失败 - CREATED：上传成功 - PUBLISHED：已发布 - TRANSCODING：转码中 - TRANSCODE_SUCCEED：转码成功 - TRANSCODE_FAILED：转码失败 - THUMBNAILING：截图中 - THUMBNAIL_SUCCEED：截图成功 - THUMBNAIL_FAILED：截图失败 - UN_REVIEW：未审核 - REVIEWING：审核中 - REVIEW_SUSPICIOUS ：审核不过，待人工复审 - REVIEW_PASSED：审核通过 - REVIEW_FAILED：审核任务失败 - REVIEW_BLOCKED：已屏蔽
	Status *[]ListAssetListRequestStatus `json:"status,omitempty"`

	// 起始时间。  格式为yyyymmddhhm mss。必须是与时区无关的UTC时间。
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间。  格式为yyyymmddhhm mss。必须是与时区无关的UTC时间。
	EndTime *string `json:"end_time,omitempty"`

	// 分类ID。
	CategoryId *int32 `json:"category_id,omitempty"`

	// 媒资标签。 单个标签不超过16个字节， 最多不超过16 个标签。 多个用英文逗号分隔，UTF8编码。
	Tags *string `json:"tags,omitempty"`

	// 在媒资标题、 描述、分类名称中模糊查询的字符串。
	QueryString *string `json:"query_string,omitempty"`

	// 音视频文件的格式，支持多格式查询，最多不超过20个。  取值如下： - 视频文件格式：MP4、TS、MOV、MXF、MPG、FLV、WMV、AVI、M4V、F4V、MPEG - 音频文件格式：MP3、OGG、WAV、WMA、APE、FLAC、AAC、AC3、MMF、AMR、M4A、M4R、WV、MP2
	MediaType *[]string `json:"media_type,omitempty"`

	// 分页编号。  默认值：0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。  取值范围：[1,100]。  默认值：10。
	Size *int32 `json:"size,omitempty"`

	// 查询顺序，按createTime顺序还是倒序
	Order *ListAssetListRequestOrder `json:"order,omitempty"`
}

func (o ListAssetListRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAssetListRequest struct{}"
	}

	return strings.Join([]string{"ListAssetListRequest", string(data)}, " ")
}

type ListAssetListRequestStatus struct {
	value string
}

type ListAssetListRequestStatusEnum struct {
	CREATING          ListAssetListRequestStatus
	FAILED            ListAssetListRequestStatus
	CREATED           ListAssetListRequestStatus
	PUBLISHED         ListAssetListRequestStatus
	WAITING_TRANSCODE ListAssetListRequestStatus
	TRANSCODING       ListAssetListRequestStatus
	TRANSCODE_SUCCEED ListAssetListRequestStatus
	TRANSCODE_FAILED  ListAssetListRequestStatus
	THUMBNAILING      ListAssetListRequestStatus
	THUMBNAIL_SUCCEED ListAssetListRequestStatus
	THUMBNAIL_FAILED  ListAssetListRequestStatus
	UN_REVIEW         ListAssetListRequestStatus
	REVIEWING         ListAssetListRequestStatus
	REVIEW_SUSPICIOUS ListAssetListRequestStatus
	REVIEW_PASSED     ListAssetListRequestStatus
	REVIEW_FAILED     ListAssetListRequestStatus
	REVIEW_BLOCKED    ListAssetListRequestStatus
}

func GetListAssetListRequestStatusEnum() ListAssetListRequestStatusEnum {
	return ListAssetListRequestStatusEnum{
		CREATING: ListAssetListRequestStatus{
			value: "CREATING",
		},
		FAILED: ListAssetListRequestStatus{
			value: "FAILED",
		},
		CREATED: ListAssetListRequestStatus{
			value: "CREATED",
		},
		PUBLISHED: ListAssetListRequestStatus{
			value: "PUBLISHED",
		},
		WAITING_TRANSCODE: ListAssetListRequestStatus{
			value: "WAITING_TRANSCODE",
		},
		TRANSCODING: ListAssetListRequestStatus{
			value: "TRANSCODING",
		},
		TRANSCODE_SUCCEED: ListAssetListRequestStatus{
			value: "TRANSCODE_SUCCEED",
		},
		TRANSCODE_FAILED: ListAssetListRequestStatus{
			value: "TRANSCODE_FAILED",
		},
		THUMBNAILING: ListAssetListRequestStatus{
			value: "THUMBNAILING",
		},
		THUMBNAIL_SUCCEED: ListAssetListRequestStatus{
			value: "THUMBNAIL_SUCCEED",
		},
		THUMBNAIL_FAILED: ListAssetListRequestStatus{
			value: "THUMBNAIL_FAILED",
		},
		UN_REVIEW: ListAssetListRequestStatus{
			value: "UN_REVIEW",
		},
		REVIEWING: ListAssetListRequestStatus{
			value: "REVIEWING",
		},
		REVIEW_SUSPICIOUS: ListAssetListRequestStatus{
			value: "REVIEW_SUSPICIOUS",
		},
		REVIEW_PASSED: ListAssetListRequestStatus{
			value: "REVIEW_PASSED",
		},
		REVIEW_FAILED: ListAssetListRequestStatus{
			value: "REVIEW_FAILED",
		},
		REVIEW_BLOCKED: ListAssetListRequestStatus{
			value: "REVIEW_BLOCKED",
		},
	}
}

func (c ListAssetListRequestStatus) Value() string {
	return c.value
}

func (c ListAssetListRequestStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListAssetListRequestStatus) UnmarshalJSON(b []byte) error {
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

type ListAssetListRequestOrder struct {
	value string
}

type ListAssetListRequestOrderEnum struct {
	ASC  ListAssetListRequestOrder
	DESC ListAssetListRequestOrder
}

func GetListAssetListRequestOrderEnum() ListAssetListRequestOrderEnum {
	return ListAssetListRequestOrderEnum{
		ASC: ListAssetListRequestOrder{
			value: "asc",
		},
		DESC: ListAssetListRequestOrder{
			value: "desc",
		},
	}
}

func (c ListAssetListRequestOrder) Value() string {
	return c.value
}

func (c ListAssetListRequestOrder) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListAssetListRequestOrder) UnmarshalJSON(b []byte) error {
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
