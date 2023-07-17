package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SubtitleModifyReq struct {

	// 媒资ID
	AssetId string `json:"asset_id"`

	// 字幕默认语言(字幕必须存在)
	DefaultLanguage *string `json:"default_language,omitempty"`

	// 需新增或修改的字幕
	AddSubtitles *[]AddSubtitle `json:"add_subtitles,omitempty"`

	// 需删除的字幕，language不能与add_subtitles重复
	DeleteSubtitles *[]DeleteSubtitle `json:"delete_subtitles,omitempty"`
}

func (o SubtitleModifyReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SubtitleModifyReq struct{}"
	}

	return strings.Join([]string{"SubtitleModifyReq", string(data)}, " ")
}
