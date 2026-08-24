# Create a ServiceStage Configuration Group and a Configuration File

This example provides best practice code for using Terraform to create a configuration group
and a configuration file which belongs to the group of ServiceStage (v3 API) in HuaweiCloud.
Both resources are metadata-type resources, creating them will not generate any fee.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* The ServiceStage service has been activated in the console (activation is free)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ServiceStage service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `configuration_group_name` - The name of the configuration group
* `configuration_name` - The name of the configuration file
* `configuration_content` - The content of the configuration file

#### Optional Variables

* `configuration_group_description` - The description of the configuration group
* `configuration_description` - The description of the configuration file

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  configuration_group_name        = "tf_test_ss_config_group"
  configuration_group_description = "Created by Terraform for ServiceStage best practice example"
  configuration_name              = "tf_test_ss_configuration"
  configuration_content           = "spring.application.name = tf-example-service"
  configuration_description       = "Created by Terraform for ServiceStage best practice example"
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
* The `name` and `description` arguments of the configuration group cannot be updated, changing
  them will recreate the group, and the same applies to the `config_group_id` and `name`
  arguments of the configuration file
* The resource name must contain `2` to `64` characters, only letters, digits, hyphens (`-`)
  and underscores (`_`) are allowed, and the name must start with a letter and end with a
  letter or a digit
* The `type` argument of the configuration file only supports `yaml` and `properties`, make
  sure the `content` argument matches the declared type, and the content supports the system
  variables of ServiceStage which are referenced by `$${VARIABLE_NAME}`
* The `content`, `type` and `description` arguments of the configuration file can be updated in
  place, each update will generate a new `version` of the configuration file
* The query detail API does not return the `type` attribute, so `lifecycle.ignore_changes`
  is used for the `type` argument in this example to avoid a permanent in-place update diff,
  if you need to change the `type` argument, remove it from `ignore_changes` first
* The `sensitive` argument is not used in this example, data encryption only takes effect for
  components deployed in containers
* Both resources support import using their IDs, note that the delete APIs always return a
  `200` status code even if the resource does not exist

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.73.7 |
