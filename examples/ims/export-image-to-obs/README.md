# Create a Export All Images To OBS Operation

This example demonstrates how to query all images owned by the current user and export them to an OBS bucket
for backup purposes in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud Provider
* HuaweiCloud access key and secret key (AK/SK)
* Appropriate IAM permissions:
  - IMS: Query images, export images
  - OBS: Create bucket, upload objects

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RAM service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `region_name` - The region where resources will be created
* `obs_bucket_name` - The name of the OBS bucket for storing exported images

#### Optional Variables

* `image_type` - The type of images to filter (ECS, BMS, etc.) (default: "")
* `image_os` - The OS of images to filter (Ubuntu, CentOS, etc.) (default: "")
* `image_name_regex` - The regex pattern to filter images by name (default: "")
* `obs_bucket_tags` - The tags of the OBS bucket (default: {})
* `file_format` - The file format of the exported image (vhd, zvhd, vmdk, raw, qcow2, zvhd2, vdi)
  (default: "zvhd2")

## Usage

* Copy this example script to your working directory.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  # OBS bucket configuration
  obs_bucket_name = "my-image-backup-bucket"

  # Image export configuration
  file_format = "zvhd2"
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

* To clean up the resources (note: this will not undo accepted/rejected invitations):

  ```bash
  $ terraform destroy
  ```

## Notes

* The export process is asynchronous and may take a long time depending on the image size
* Only active images are exported (images with status "active")
* The configuration uses `count` instead of `for_each` to iterate over active images
* A `locals` block is used to filter active images, making the code more maintainable
* The default timeout for image export is 20 minutes per image
* Make sure the OBS bucket has sufficient storage capacity
* The bucket URL format is: `bucket-name:image-name-image-id.file-format`

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.69.0 |
