# Create a VOD Media Category and a Media Asset by URL

This example provides best practice code for using Terraform to create a media category and
a media asset which is uploaded by URL of Video on Demand (VOD) in HuaweiCloud. The URL upload
mode pulls the media source file asynchronously, no local media file or OBS bucket is required.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* The VOD service has been activated in the console (activation is free)
* An accessible HTTP or HTTPS URL of a media source file (referenced by the `media_asset_url`
  variable)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the VOD service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `media_category_name` - The name of the media category
* `media_asset_name` - The name of the media asset
* `media_asset_url` - The HTTP or HTTPS URL of the media source file

#### Optional Variables

* `media_asset_description` - The description of the media asset
* `media_asset_labels` - The labels of the media asset, separated by commas

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  media_category_name     = "tf_test_vod_asset_category"
  media_asset_name        = "tf_test_vod_media_asset"
  media_asset_url         = "https://test-videos.co.uk/vids/bigbuckbunny/mp4/h264/360/Big_Buck_Bunny_360_10s_1MB.mp4"
  media_asset_description = "Created by Terraform for VOD best practice example"
  media_asset_labels      = "tf_label_1,tf_label_2"
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
* The `media_asset_url` argument only supports the `http` and `https` protocols, the URL upload
  pulls the media source file asynchronously, the creation will succeed even if the URL is
  temporarily inaccessible, but the media asset will remain in a failed upload status in the
  console, replace it with your own accessible URL
* The `media_type` argument must match the actual format of the media source file, and it
  cannot be updated, changing it will recreate the media asset
* The `category_id` argument of the media asset references the media category created in this
  example, if it is not specified, the media asset will be classified into the system preset
  `Other` category
* The `publish` argument is not used in this example and defaults to `false`, the media asset
  will not be published to the CDN, so no traffic fee is generated, and the `template_group_name`
  argument is not used either, so no transcoding task will be triggered
* The source file is copied to the VOD bucket by default (`storage_mode` is `0`), a small amount
  of storage fee may be generated, remember to run `terraform destroy` to delete the media asset
  after testing
* Both resources support import, the import ID of the media category is its numeric ID, and the
  import ID of the media asset is the asset ID, note that the `url` argument of the media asset
  will not be read back after import

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.38.0 |
