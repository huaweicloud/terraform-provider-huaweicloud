# Create a transcoding template
resource "huaweicloud_mpc_transcoding_template" "test" {
  name                  = var.template_name
  low_bitrate_hd        = var.low_bitrate_hd
  dash_segment_duration = var.dash_segment_duration
  hls_segment_duration  = var.hls_segment_duration
  output_format         = var.output_format

  audio {
    bitrate       = var.audio_bitrate
    channels      = var.audio_channels
    codec         = var.audio_codec
    output_policy = var.audio_output_policy
    sample_rate   = var.audio_sample_rate
  }

  video {
    max_consecutive_bframes = var.video_max_consecutive_bframes
    bitrate                 = var.video_bitrate
    black_bar_removal       = var.video_black_bar_removal
    codec                   = var.video_codec
    fps                     = var.video_fps
    level                   = var.video_level
    max_iframes_interval    = var.video_max_iframes_interval
    output_policy           = var.video_output_policy
    quality                 = var.video_quality
    profile                 = var.video_profile
    height                  = var.video_height
    width                   = var.video_width
  }
}
