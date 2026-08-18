# Create a VOD Media Category and a Transcoding Template Group

This example provides best practice code for using Terraform to create a media category and
a transcoding template group of Video on Demand (VOD) in HuaweiCloud. Both resources are
configuration-type resources, creating them will not generate any fee.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* The VOD service has been activated in the console (activation is free)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the VOD service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `media_category_name` - The name of the media category
* `template_group_name` - The name of the transcoding template group

#### Optional Variables

* `template_group_description` - The description of the transcoding template group

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  media_category_name        = "tf_test_vod_media_category"
  template_group_name        = "tf_test_vod_template_group"
  template_group_description = "Created by Terraform for VOD best practice example"
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

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The `parent_id` argument of the media category cannot be updated, changing it will recreate
  the category, the parent category ID of a top-level category is `0` by default, and a media
  category that still contains child categories or media assets cannot be deleted
* The `quality_info` block of the transcoding template group is required, the valid values of
  `output_format` are `MP4`, `DASH`, `DASH_HLS`, `MP3`, `ADTS` and `HLS`, and the valid values
  of the video `quality` are `4K`, `2K`, `FHD`, `HD`, `SD` and `LD`
* The `sample_rate` argument of the audio is an enumerated number, `1` means `AUTO`, `2` means
  `22,050` Hz, `3` means `32,000` Hz, `4` means `44,100` Hz, `5` means `48,000` Hz and `6`
  means `96,000` Hz
* If the output format is `HLS`, the `hls_segment_duration` argument can be set to specify the
  segment duration, valid values are `2`, `3`, `5` and `10`
* The `watermark_template_ids` argument requires existing watermark templates, and creating a
  watermark template needs a local watermark image file, so this example does not use it
* Both resources support import, the import ID of the media category is its numeric ID, and the
  import ID of the transcoding template group is the group ID, note that the `parent_id`
  argument of the media category will not be read back after import

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.37.0 |
