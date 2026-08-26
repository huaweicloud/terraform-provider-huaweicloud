# Associate an EIP to a ServiceStage Environment

This example provides best practice code for using Terraform to create a VPC, an EIP, a
ServiceStage environment (v3 API) bound to the VPC, and associate the EIP to the environment
in HuaweiCloud.

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
* `eip_bandwidth_name` - The name of the dedicated bandwidth of the EIP
* `environment_name` - The name of the environment

#### Optional Variables

* `environment_description` - The description of the environment

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                = "tf_test_ss_vpc"
  eip_bandwidth_name      = "tf_test_ss_eip_bandwidth"
  environment_name        = "tf_test_ss_environment"
  environment_description = "Created by Terraform for ServiceStage best practice example"
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

* The EIP created in this example will incur fees (pay-per-use, billed by traffic), please
  run `terraform destroy` in time after the verification is completed
* The name of the environment must contain `2` to `64` characters, only letters, digits,
  hyphens (`-`) and underscores (`_`) are allowed, and the name must start with a letter and
  end with a letter or a digit, the description must not exceed `128` characters
* The `vpc_id` argument of the environment cannot be updated, changing it will recreate the
  environment, and the `environment_id` argument of the associate resource is also `ForceNew`
* The creation, update and deletion of the associate resource all use the PUT method which
  manages the associated resources in full-set mode, do not associate other resources to the
  same environment in the console at the same time, otherwise they will be overwritten or
  dissociated by Terraform
* The `name` argument of the `resources` block is not set in this example because the server
  does not return the name for EIP-type resources, if you associate CCE clusters instead,
  the cluster name must be configured, note that the `name` argument is only supported by
  the provider of `1.93.0` and later, in earlier versions the `resources` block is an
  unordered set which only supports the `id` and `type` arguments
* Destroying the associate resource only dissociates the resources from the environment, the
  EIP itself will be deleted by its own resource definition
* The associate resource supports import using the environment ID

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.69.0 |
