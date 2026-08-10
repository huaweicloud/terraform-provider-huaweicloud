# Create a Template with a New Version

This example provides best practice code for using Terraform to create an RFS template in HuaweiCloud Resource
Formation Service (RFS). The example first creates a template with an inline Terraform template body, and then
publishes a new version for the template with an extended template body.

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

* `template_name` - The name of the RFS template
* `template_body` - The Terraform template body of the RFS template
* `template_version_body` - The Terraform template body of the RFS template version

#### Optional Variables

* `template_description` - The description of the RFS template (default: `""`)
* `template_initial_version_description` - The initial version description of the RFS template (default: `""`)
* `template_version_description` - The description of the RFS template version (default: `""`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  template_name         = "tf-test-template"
  template_body         = <<-EOF
  resource "huaweicloud_vpc" "test" {
    name = "tf-test-vpc"
    cidr = "172.16.0.0/16"
  }
  EOF

  template_version_body = <<-EOF
  resource "huaweicloud_vpc" "test" {
    name = "tf-test-vpc"
    cidr = "172.16.0.0/16"
  }

  resource "huaweicloud_vpc_subnet" "test" {
    name       = "tf-test-subnet"
    vpc_id     = huaweicloud_vpc.test.id
    cidr       = "172.16.1.0/24"
    gateway_ip = "172.16.1.1"
  }
  EOF
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
* The `template_name` only allows letters, digits, hyphens (-), and underscores (_), and must be unique within the
  region
* The `template_body` must be a valid Terraform configuration in HCL format, and only the template definition is
  stored, the resources in the template body will not be created by this example
* All arguments of the template and the template version (except `template_description`) are non-updatable, changing
  them will recreate the resource
* Deleting the template will also delete all its versions

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.91.0 |
