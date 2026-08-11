# Create a Resource Stack with an Execution Plan

This example provides best practice code for using Terraform to create an RFS resource stack and an execution plan
in HuaweiCloud Resource Formation Service (RFS). The example first creates an empty resource stack, and then creates
an execution plan with an inline Terraform template body to preview the resource changes before deployment.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RFS resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `stack_name` - The name of the RFS resource stack
* `execution_plan_name` - The name of the RFS execution plan
* `execution_plan_template_body` - The Terraform template body used by the RFS execution plan

#### Optional Variables

* `stack_description` - The description of the RFS resource stack (default: `""`)
* `execution_plan_description` - The description of the RFS execution plan (default: `""`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  stack_name          = "tf-test-stack"
  execution_plan_name = "tf-test-plan"

  execution_plan_template_body = <<-EOT
  terraform {
    required_providers {
      huaweicloud = {
        source  = "huawei.com/provider/huaweicloud"
        version = ">= 1.41.0"
      }
    }
  }

  resource "huaweicloud_vpc" "test" {
    name = "tf-test-vpc"
    cidr = "192.168.0.0/16"
  }

  resource "huaweicloud_vpc_subnet" "test" {
    vpc_id     = huaweicloud_vpc.test.id
    name       = "tf-test-subnet"
    cidr       = "192.168.1.0/24"
    gateway_ip = "192.168.1.1"
  }
  EOT
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
* This example creates an empty resource stack without a template, so no cloud resources will be deployed, and no
  agency authorization is required
* The execution plan only previews the resource changes described by the template body, it does not actually create
  any cloud resources
* The template body must declare the `required_providers` block with the source `huawei.com/provider/huaweicloud`,
  otherwise the RFS execution plan engine cannot resolve the provider packages, and the plan creation will fail
* All arguments of the execution plan are non-updatable, changing them will recreate the resource
* Deleting the resource stack will also delete the related execution plans

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.91.0 |
