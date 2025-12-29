# Create a SMS migration project example

This example provides best practice code for using Terraform to create a migration project in HuaweiCloud SMS service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `name` - The migration project name.
* `region` - The region name.
* `use_public_ip` - Whether to use a public IP address for migration.
* `exist_server` - Whether the server already exists.
* `type` - The migration project type.
* `syncing` - whether to continue syncing after the first copy or sync.

#### Optional Variables

* `description` - The migration project description.
* `is_default` - Whether to use the default template (default: false).
* `start_target_server` - Whether to start the destination virtual machine after migration (default: false).
* `speed_limit` - The migration rate limit in Mbps (default: "0").
* `enterprise_project` - The name of the enterprise project (default: "default").
* `start_network_check` - Whether to enable network quality detection (default: false).

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  migration_project_name          = "tf_test_sms_migration_project"
  migration_project_region        = "tf_test_sms_migration_project"
  migration_project_use_public_ip = true
  migration_project_exist_server  = true
  migration_project_type          = "tf_test_sms_migration_project"
  migration_project_syncing       = true
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

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.75.1 |
