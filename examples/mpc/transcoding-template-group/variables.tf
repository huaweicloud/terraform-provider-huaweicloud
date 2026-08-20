# Variable definitions for authentication
variable "region_name" {
  description = "The region where the MPC resources are located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
  sensitive   = true
}

# Variable definitions for huaweicloud_mpc_transcoding_template_group
variable "template_name" {
  description = "The name of the transcoding template group"
  type        = string
}

variable "low_bitrate_hd" {
  description = "Whether to enable low bitrate HD. true: enable, false: disable"
  type        = bool
  default     = true
}

variable "dash_segment_duration" {
  description = "The DASH segment duration in seconds"
  type        = number
  default     = 5
}

variable "hls_segment_duration" {
  description = "The HLS segment duration in seconds"
  type        = number
  default     = 5
}

variable "output_format" {
  description = "The output format. 1: HLS, 2: DASH, 3: HLS+DASH, 4: MP4, 5: MP3, 6: ADTS"
  type        = number
  default     = 1
}

# Audio parameters
variable "audio_bitrate" {
  description = "The audio bitrate. 0: auto"
  type        = number
  default     = 0
}

variable "audio_channels" {
  description = "The audio channels. 1: AUTO, 2: mono, 6: stereo"
  type        = number
  default     = 2
}

variable "audio_codec" {
  description = "The audio codec. 1: AAC, 2: HEAAC1, 3: HEAAC2, 4: MP3"
  type        = number
  default     = 2
}

variable "audio_output_policy" {
  description = "The audio output policy. transcode: transcoding, copy: passthrough, discard: discard"
  type        = string
  default     = "transcode"
}

variable "audio_sample_rate" {
  description = "The audio sample rate. 1: AUTO, 2: 22050Hz, 3: 32000Hz, 4: 44100Hz, 5: 48000Hz, 6: 96000Hz"
  type        = number
  default     = 1
}

# Video common parameters
variable "video_max_consecutive_bframes" {
  description = "The maximum number of consecutive B-frames"
  type        = number
  default     = 7
}

variable "video_black_bar_removal" {
  description = "Whether to remove black bars. 0: disable, 1: enable"
  type        = number
  default     = 0
}

variable "video_codec" {
  description = "The video codec. 1: H.264, 2: H.265"
  type        = number
  default     = 2
}

variable "video_fps" {
  description = "The video frame rate. 0: auto"
  type        = number
  default     = 0
}

variable "video_level" {
  description = "The video level. 15: default"
  type        = number
  default     = 15
}

variable "video_max_iframes_interval" {
  description = "The maximum interval between I-frames"
  type        = number
  default     = 5
}

variable "video_output_policy" {
  description = "The video output policy. transcode: transcoding, copy: passthrough, discard: discard"
  type        = string
  default     = "transcode"
}

variable "video_quality" {
  description = "The video quality. 1: VBR, 2: CBR"
  type        = number
  default     = 1
}

variable "video_profile" {
  description = "The video profile. 1: baseline, 2: main, 3: high, 4: default"
  type        = number
  default     = 4
}

# Video output definitions
variable "video_outputs" {
  description = "The list of video output definitions"
  type        = list(object({
    width   = number
    height  = number
    bitrate = number
  }))
  default     = [
    {
      width   = 1920
      height  = 1080
      bitrate = 0
    }
  ]
  nullable    = false
}
