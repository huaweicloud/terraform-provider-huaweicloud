package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateThumbReq struct {
	Input *ObsObjInfo `json:"input"`

	Output *ObsObjInfo `json:"output"`

	// 用户自定义数据。
	UserData *string `json:"user_data,omitempty"`

	ThumbnailPara *ThumbnailPara `json:"thumbnail_para"`

	// 是否压缩抽帧图片生成tar包。  取值如下： - 0：压缩。 - 1：不压缩 默认值：1
	Tar *int32 `json:"tar,omitempty"`

	// 是否同步处理，同步处理是指不下载全部文件，快速定位到截图位置进行截图。  取值如下： - 0：排队处理。 - 1：同步处理，暂只支持按时间点截单张图。 默认值：0
	Sync *int32 `json:"sync,omitempty"`

	// 是否使用原始输出目录。  取值如下： - 0：不使用原始输出目录，下发的输出目录后面追加随机目录，防止截图文件outputUri相同被覆盖。 - 1：使用原始输出目录。 默认值：0
	OriginalDir *int32 `json:"original_dir,omitempty"`
}

func (o CreateThumbReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateThumbReq struct{}"
	}

	return strings.Join([]string{"CreateThumbReq", string(data)}, " ")
}
