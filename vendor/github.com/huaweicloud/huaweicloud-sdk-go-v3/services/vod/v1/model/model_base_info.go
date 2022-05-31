package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 媒资基本信息。
type BaseInfo struct {

	// 媒资标题。  长度不超过128个字节，UTF8编码。
	Title *string `json:"title,omitempty"`

	// 媒资文件名。
	VideoName *string `json:"video_name,omitempty"`

	// 媒资描述。  长度不超过1024个字节。
	Description *string `json:"description,omitempty"`

	// 媒资分类id。
	CategoryId *int64 `json:"category_id,omitempty"`

	// 媒资分类名称。
	CategoryName *string `json:"category_name,omitempty"`

	// 媒资创建时间。  格式为yyyymmddhhmmss。必须是与时区无关的UTC时间。
	CreateTime *string `json:"create_time,omitempty"`

	// 媒资最近修改时间。  格式为yyyymmddhhmmss。必须是与时区无关的UTC时间。
	LastModified *string `json:"last_modified,omitempty"`

	// 音视频文件类型。  取值如下： - 视频文件：MP4、TS、MOV、MXF、MPG、FLV、WMV、AVI、M4V、F4V、MPEG、3GP、ASF、MKV。 - 音频文件：MP3、OGG、WAV、WMA、APE、FLAC、AAC、AC3、MMF、AMR、M4A、M4R、WV、MP2。
	VideoType *string `json:"video_type,omitempty"`

	// 媒资标签。  单个标签不超过16个字节，最多不超过16个标签。  多个用逗号分隔，UTF8编码。
	Tags *string `json:"tags,omitempty"`

	MetaData *MetaData `json:"meta_data,omitempty"`

	// 原始视频文件的访问地址。
	VideoUrl *string `json:"video_url,omitempty"`

	// 原视频文件的OBS临时访问地址,仅媒资详情接口生效
	SignUrl *string `json:"sign_url,omitempty"`

	// 封面信息。
	CoverInfoArray *[]CoverInfo `json:"cover_info_array,omitempty"`

	// 字幕信息数组
	SubtitleInfo *[]SubtitleInfo `json:"subtitle_info,omitempty"`

	SourcePath *FileAddr `json:"source_path,omitempty"`

	OutputPath *FileAddr `json:"output_path,omitempty"`
}

func (o BaseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BaseInfo struct{}"
	}

	return strings.Join([]string{"BaseInfo", string(data)}, " ")
}
