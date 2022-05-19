package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type QualityEnhanceVideo struct {
	VideoDenoise *VideoDenoise `json:"video_denoise,omitempty"`

	VideoSharp *VideoSharp `json:"video_sharp,omitempty"`

	VideoContrast *VideoContrast `json:"video_contrast,omitempty"`

	VideoSuperresolution *VideoSuperresolution `json:"video_superresolution,omitempty"`

	VideoDeblock *VideoDeblock `json:"video_deblock,omitempty"`

	VideoSaturation *VideoSaturation `json:"video_saturation,omitempty"`
}

func (o QualityEnhanceVideo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QualityEnhanceVideo struct{}"
	}

	return strings.Join([]string{"QualityEnhanceVideo", string(data)}, " ")
}
