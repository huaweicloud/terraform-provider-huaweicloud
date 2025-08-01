# Create a Workspace APP policy group with scaling policy

This example provides best practice code for using Terraform to create a Workspace APP policy group with scaling policy
in HuaweiCloud Workspace service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* Workspace service enabled in the target region

## Required Variables

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
* `app_group_name` - The name of the APP group
* `policy_group_name` - The name of the policy group
* `max_scaling_amount` - The maximum number of instances that can be scaled out
* `single_expansion_count` - The number of instances to scale out in a single scaling operation

#### Optional Variables

* `app_server_group_app_type` - The application type of the APP server group (default: "SESSION_DESKTOP_APP")
* `app_server_group_os_type` - The operating system type of the APP server group (default: "Windows")
* `app_server_group_system_disk_type` - The system disk type of the APP server group (default: "SAS")
* `app_server_group_system_disk_size` - The system disk size of the APP server group in GB (default: 80)
* `policy_group_priority` - The priority of the policy group (default: 1)
* `policy_group_description` - The description of the policy group (default: "Created APP policy group by Terraform")
* `target_type` - The type of target for the policy group (default: "APPGROUP")
* `automatic_reconnection_interval` - The automatic reconnection interval in minutes (default: 10)
* `session_persistence_time` - The session persistence time in minutes (default: 120)
* `forbid_screen_capture` - Whether to forbid screen capture (default: true)
* `session_usage_threshold` - The session usage threshold percentage (default: 80)
* `shrink_after_session_idle_minutes` - The number of minutes to wait before shrinking idle instances (default: 30)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  app_server_group_name             = "tf_test_server_group"
  app_server_group_flavor_id        = "workspace.appstream.general.xlarge.4"
  app_server_group_image_id         = "2ac7b1fb-b198-422b-a45f-61ea285cb6e7"
  app_server_group_image_product_id = "OFFI886188719633408000"
  app_group_name                    = "tf_test_app_group"
  policy_group_name                 = "tf_test_policy_group"
  max_scaling_amount                = 10
  single_expansion_count            = 2
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
* The scaling policy is dependent on the APP server group, APP group, and policy group
* Please read the implicit and explicit dependencies in the script carefully
* All resources will be created in the specified region
* The scaling policy supports session-based scaling configuration

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.77.1 |
