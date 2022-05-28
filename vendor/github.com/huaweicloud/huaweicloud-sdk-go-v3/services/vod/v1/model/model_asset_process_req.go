package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AssetProcessReq struct {

	// 媒资ID。
	AssetId string `json:"asset_id"`

	// 转码模板组名称。   若不为空，则使用指定的转码模板对上传的音视频进行转码，您可以在视频点播控制台配置转码模板，具体请参见[转码设置](https://support.huaweicloud.com/usermanual-vod/vod_01_0072.html)。
	TemplateGroupName *string `json:"template_group_name,omitempty"`

	// 是否自动加密。  取值如下： - 0：表示不加密。 - 1：表示需要加密。  默认值：0。  加密与转码必须要一起进行，当需要加密时，转码参数不能为空，且转码输出格式必须要为HLS。
	AutoEncrypt *int32 `json:"auto_encrypt,omitempty"`

	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`

	// 字幕文件ID。  > 仅在[创建媒资](https://support.huaweicloud.com/api-vod/vod_04_0196.html)时，请求参数设置了“**subtitles**”时，该参数设置才生效。
	SubtitleId *[]int32 `json:"subtitle_id,omitempty"`
}

func (o AssetProcessReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssetProcessReq struct{}"
	}

	return strings.Join([]string{"AssetProcessReq", string(data)}, " ")
}
