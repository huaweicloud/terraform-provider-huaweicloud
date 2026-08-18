resource "huaweicloud_vod_media_category" "test" {
  name = var.media_category_name
}

resource "huaweicloud_vod_transcoding_template_group" "test" {
  name        = var.template_group_name
  description = var.template_group_description
  audio_codec = "HEAAC1"
  video_codec = "H265"

  quality_info {
    output_format = "MP4"

    audio {
      channels    = 1
      sample_rate = 2
    }

    video {
      bitrate    = 1000
      frame_rate = 1
      height     = 1080
      quality    = "FHD"
      width      = 1920
    }
  }
}
