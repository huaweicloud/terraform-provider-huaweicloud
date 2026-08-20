# Create a transcoding template group

This example provides best practice code for using Terraform to create a transcoding template group in HuaweiCloud MPC
(Media Processing Center) service, which defines common audio and video parameters along with multiple output qualities.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the MPC resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `template_name` - The name of the transcoding template group

#### Optional Variables

* `output_format` - The output format. 1: HLS, 2: DASH, 3: HLS+DASH, 4: MP4, 5: MP3, 6: ADTS (default: 1)
* `low_bitrate_hd` - Whether to enable low bitrate HD (default: true)
* `hls_segment_duration` - The HLS segment duration in seconds (default: 5)
* `dash_segment_duration` - The DASH segment duration in seconds (default: 5)
* `audio_codec` - The audio codec. 1: AAC, 2: HEAAC1, 3: HEAAC2, 4: MP3 (default: 2)
* `audio_sample_rate` - The audio sample rate .
  Valid values: 1=AUTO, 2=22050Hz, 3=32000Hz, 4=44100Hz, 5=48000Hz, 6=96000Hz (default: 1)
* `audio_channels` - The audio channels. 1: AUTO, 2: mono, 6: stereo (default: 2)
* `audio_bitrate` - The audio bitrate. 0: auto (default: 0)
* `audio_output_policy` - The audio output policy. transcode, copy, or discard (default: "transcode")
* `video_codec` - The video codec. 1: H.264, 2: H.265 (default: 2)
* `video_profile` - The video profile. 1: baseline, 2: main, 3: high, 4: default (default: 4)
* `video_level` - The video level. 15: default (default: 15)
* `video_quality` - The video quality. 1: VBR, 2: CBR (default: 1)
* `video_fps` - The video frame rate. 0: auto (default: 0)
* `video_max_iframes_interval` - The maximum interval between I-frames (default: 5)
* `video_max_consecutive_bframes` - The maximum number of consecutive B-frames (default: 7)
* `video_black_bar_removal` - Whether to remove black bars. 0: disable, 1: enable (default: 0)
* `video_output_policy` - The video output policy. transcode, copy, or discard (default: "transcode")
* `video_outputs` - The list of video output definitions, each containing width, height, and bitrate.
  Default: 1920x1080, bitrate=0

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  template_name = "your_transcoding_template_group_name"
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The transcoding template group uses `video_common` for shared video parameters and `videos` for per-output settings
* Add multiple entries to `video_outputs` to define multiple output qualities in a single template group
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.37.0 |
