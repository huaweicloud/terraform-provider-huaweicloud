# Create a ai-ops setting

This example provides best practice code for using Terraform to create a ai-ops setting in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CSS cluster is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `cluster_name` - The CSS cluster name

#### Optional Variables

* `availability_zone` - The availability zone to which the CSS cluster belongs (default: "")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `cluster_flavor` - The flavor of the CSS cluster (default: "")
* `cluster_engine_version` - The engine version of the CSS cluster (default: "7.10.2")
* `cluster_instance_number` - The number of instances of the CSS cluster (default: 3)
* `cluster_volume_type` - The volume type of the CSS cluster (default: "ULTRAHIGH")
* `cluster_volume_size` - The volume size of the CSS cluster (default: 40)
* `ai_ops_check_type` - The check type of the AI Ops setting (default: "full_detection")
* `ai_ops_period` - The period of the AI Ops setting (default: "12:00 GMT+08:00")
* `ai_ops_check_items` - The check items of the AI Ops setting, only required when `ai_ops_check_type`
  is `partial_detection` (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  cluster_name        = "your_cluster_name"
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
* The creation of the CSS cluster takes about 30 to 60 minutes depending on the flavor and instance number
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.88.0 |
