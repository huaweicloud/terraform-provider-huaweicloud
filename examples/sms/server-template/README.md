# Create a SMS server template example

This example provides best practice code for using Terraform to create a server template in HuaweiCloud SMS service.

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

* `name` - The server template name.
* `availability_zone` - The availability zone where the target server is located.

#### Optional Variables

* `region` - The region where the target server is located.
* `project_id` - The project ID where the target server is located.
* `vpc_id` - The ID of the VPC which the target server belongs to.
* `subnet_ids` - The list of subnet IDs to attach to the target server.
* `security_group_ids` - The list of security group IDs to associate with the target server.
* `volume_type` - The disk type of the target server.
* `flavor` - The flavor ID for the target server.
* `target_server_name` - The name of the target server.
* `bandwidth_size` - The bandwidth size in Mbit/s about the public IP address that will be used for migration.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  server_template_name = "tf_test_sms_server_template"
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
| huaweicloud | >= >= 1.37.0 |
