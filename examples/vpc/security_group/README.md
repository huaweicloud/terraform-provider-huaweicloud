# Create a basic security group and rules

This example provides best practice code for using Terraform to create a basic security group and its rules in
HuaweiCloud VPC service. The security group is used to control network traffic using flexible rule configurations.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the security group is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `security_group_name` - The name of the security group

#### Optional Variables

* `security_group_rule_configurations` - The list of security group rule configurations
  - `direction` - The direction of the security group rule
  - `ethertype` - The ethertype of the security group rule
  - `protocol` - The protocol of the security group rule
  - `ports` - The port range of the security group rule
  - `remote_ip_prefix` - The remote IP prefix of the security group rule

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  region_name         = "your_region"
  access_key          = "your_access_key"
  secret_key          = "your_secret_key"
  security_group_name = "tf_test_secgroup"
  # security_group_rule_configurations is optional, no rules will be kept
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
* This example creates a security group with configurable rules
* Default rules allow ICMP ingress and TCP port 443 ingress if no custom rules are specified
* All resources will be created in the specified region
* The security group can be associated with other resources like ECS instances

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.56.0 |
