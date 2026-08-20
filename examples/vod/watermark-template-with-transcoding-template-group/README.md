# Create a VOD Watermark Template and a Transcoding Template Group

This example provides best practice code for using Terraform to create a watermark template
and a transcoding template group which references the watermark template of Video on Demand
(VOD) in HuaweiCloud. Both resources are configuration-type resources, creating them will not
generate any fee.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* The VOD service has been activated in the console (activation is free)
* A local watermark image file in `PNG` format (referenced by the `watermark_template_image_file`
  variable)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the VOD service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `watermark_template_name` - The name of the watermark template
* `watermark_template_image_file` - The absolute path of the local watermark image file in PNG format
* `template_group_name` - The name of the transcoding template group

#### Optional Variables

* `template_group_description` - The description of the transcoding template group

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  watermark_template_name       = "tf_test_vod_watermark_template"
  watermark_template_image_file = "D:/styles/watermark.png"
  template_group_name           = "tf_test_vod_watermark_group"
  template_group_description    = "Created by Terraform for VOD watermark example"
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
* The `image_file` argument is the local file path on the machine running Terraform, the
  provider reads this file and uploads it to the upload URL returned by the VOD API, the
  creation will fail if the file does not exist
* The `image_file` and `image_type` arguments of the watermark template cannot be updated,
  changing them will recreate the template, and the `image_file` argument will not be read
  back after import
* The creation of the watermark template includes three steps (creating the template,
  uploading the image and confirming the upload), if the image upload fails, an incomplete
  template may remain in the console and needs to be deleted manually
* The `width`, `height`, `horizontal_offset` and `vertical_offset` arguments are strings
  representing ratios or relative offsets, e.g. `0.01`, do not pass numeric values
* The `quality_info` block of the transcoding template group is required, the valid values of
  `output_format` are `MP4`, `DASH`, `DASH_HLS`, `MP3`, `ADTS` and `HLS`, and the valid values
  of the video `quality` are `4K`, `2K`, `FHD`, `HD`, `SD` and `LD`
* The transcoding template group references the watermark template through
  `watermark_template_ids`, so `destroy` will delete the template group first automatically
* Both resources support import, the import ID of the watermark template is its template ID,
  and the import ID of the transcoding template group is the group ID

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.37.0 |
