package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateAssetResponse struct {

	// 媒资ID。
	AssetId *string `json:"asset_id,omitempty"`

	// 视频上传地址。  用于客户端上传不超过20MB的小视频文件（<=20MB）。该URL中携带了临时授权信息，当文件大于20MB时，需要采用分段方式上传。  > 您可以使用PUT请求向“**video_upload_url**”中上传视频文件。其中，“**Content-Type**”需要根据上传的视频文件类型进行设置，如下所示：视频文件：video/视频格式，如：“video/mp4”。
	VideoUploadUrl *string `json:"video_upload_url,omitempty"`

	// 封面上传地址，用于上传封面文件。  > 您可以使用PUT请求向“**cover_upload_url**”中上传封面图片。其中，“**Content-Type**”需要根据上传的封面文件类型进行设置，如下所示：图片文件：image/图片格式，如：“image/png”。
	CoverUploadUrl *string `json:"cover_upload_url,omitempty"`

	// 字幕上传地址，用于上传字幕。  > 您可以使用PUT请求向“**subtitle_upload_urls**”中上传字幕文件。其中，“**Content-Type**”需要根据上传的字幕文件类型进行设置，如下所示：字幕文件：application/octet-stream。
	SubtitleUploadUrls *[]string `json:"subtitle_upload_urls,omitempty"`
	HttpStatusCode     int       `json:"-"`
}

func (o UpdateAssetResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAssetResponse struct{}"
	}

	return strings.Join([]string{"UpdateAssetResponse", string(data)}, " ")
}
