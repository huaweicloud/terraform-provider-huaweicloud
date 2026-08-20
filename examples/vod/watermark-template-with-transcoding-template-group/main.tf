resource "huaweicloud_vod_watermark_template" "test" {
  name       = var.watermark_template_name
  image_file = var.watermark_template_image_file
  image_type = "PNG"
}

resource "huaweicloud_vod_transcoding_template_group" "test" {
  name                   = var.template_group_name
  description            = var.template_group_description
  audio_codec            = "HEAAC1"
  video_codec            = "H265"
  watermark_template_ids = [huaweicloud_vod_watermark_template.test.id]

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
