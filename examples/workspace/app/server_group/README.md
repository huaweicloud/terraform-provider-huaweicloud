# Create a Workspace APP server group

This example provides best practice code for using Terraform to create a Workspace APP server group in HuaweiCloud
Workspace service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* Workspace service enabled in the target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the Workspace service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `app_server_group_name` - The name of the APP server group
* `app_server_group_flavor_id` - The flavor ID of the APP server group
* `app_server_group_image_id` - The image ID of the APP server group
* `app_server_group_image_product_id` - The image product ID of the APP server group

#### Optional Variables

* `app_server_group_app_type` - The application type of the APP server group (default: "SESSION_DESKTOP_APP")
* `app_server_group_os_type` - The operating system type of the APP server group (default: "Windows")
* `app_server_group_system_disk_type` - The system disk type of the APP server group (default: "SAS")
* `app_server_group_system_disk_size` - The system disk size of the APP server group in GB (default: 80)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  app_server_group_name             = "tf_test_server_group"
  app_server_group_flavor_id        = "workspace.appstream.general.xlarge.4"
  app_server_group_image_id         = "2ac7b1fb-b198-422b-a45f-61ea285cb6e7"
  app_server_group_image_product_id = "OFFI886188719633408000"
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
* The server group is dependent on the Workspace service configuration
* Please read the implicit and explicit dependencies in the script carefully
* All resources will be created in the specified region
* The server group supports different application types and operating systems

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.72.0 |
