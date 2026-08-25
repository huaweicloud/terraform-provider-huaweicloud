# Create a ServiceStage Environment, an Application and Its Application Configuration

This example provides best practice code for using Terraform to create a VPC, a ServiceStage
environment (v3 API) bound to the VPC, an application, and the application-level environment
variable configuration of the application in HuaweiCloud.
The VPC and all ServiceStage resources in this example are free of charge.

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

* `vpc_name` - The name of the VPC
* `environment_name` - The name of the environment
* `application_name` - The name of the application
* `application_configuration_env_name` - The name of the environment variable
* `application_configuration_env_value` - The value of the environment variable

#### Optional Variables

* `environment_description` - The description of the environment
* `application_description` - The description of the application

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                            = "tf_test_ss_vpc"
  environment_name                    = "tf_test_ss_environment"
  environment_description             = "Created by Terraform for ServiceStage best practice example"
  application_name                    = "tf_test_ss_application"
  application_description             = "Created by Terraform for ServiceStage best practice example"
  application_configuration_env_name  = "TF_EXAMPLE_APP_ENV"
  application_configuration_env_value = "tf_example_value"
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

* The name of the environment and the application must contain `2` to `64` characters, only
  letters, digits, hyphens (`-`) and underscores (`_`) are allowed, and the name must start
  with a letter and end with a letter or a digit, the description must not exceed `128`
  characters
* The `vpc_id` argument of the environment cannot be updated, changing it will recreate the
  environment, and the `deploy_mode` argument is not set in this example, the default value
  returned by the server will be used
* The name of the environment variable must contain `1` to `64` characters and must start with
  a letter, a hyphen (`-`) or an underscore (`_`)
* The creation and update of the application configuration both use the PUT method which
  overwrites all environment variables of the application under the specified environment,
  do not manage the same application and environment combination in other places at the same
  time
* The `environment_id` and `application_id` arguments of the application configuration cannot
  be updated, changing them will recreate the resource
* The `assign_strategy` argument is not set in this example, the effective strategy defaults
  to the continuously effective, if it is set to `true`, the environment variables only take
  effect when the component is first created
* The `enterprise_project_id` argument is not set in this example, the environment and the
  application will be created in the default enterprise project
* The environment and the application support import using their IDs, the application
  configuration supports import using the format of `<environment_id>/<application_id>`

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.73.9 |
