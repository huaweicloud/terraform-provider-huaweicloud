# Create a high availability and secure GES graph instance

This example provides best practice code for using Terraform to create a cross-AZ GES graph instance with
private network access and security hardening in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* GES service enabled in the target region
* At least two availability zones in the target region (required for cross-AZ deployment)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GES resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `graph_name` - The name of the GES graph instance

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `graph_size_type_index` - The graph size type index (default: "3", 100 million edges)
* `enable_multi_az` - Whether to enable cross-AZ deployment (default: true)
* `enable_https` - Whether to enable HTTPS security mode (default: false)
* `enable_audit` - Whether to enable graph audit logs via LTS (default: false)
* `audit_log_group_name` - The LTS log group name for audit logs (default: "", required when `enable_audit` is true)
* `graph_tags` - The tags of the GES graph instance (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  graph_name          = "your_graph_name"
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

## Features

This example demonstrates the following features:

1. **Cross-AZ High Availability**: Deploys graph ECS instances across two AZs via `enable_multi_az`
2. **Private Network Access**: Uses VPC and subnet for private network connectivity
3. **Security Group Restriction**: Restricts inbound traffic to the VPC CIDR only
4. **LTS Audit Logging**: Optional graph operation audit via LTS

## Network Access Best Practice

To prevent service interruption caused by floating IP switchover, applications should poll the physical IP
addresses in the graph instance `traffic_ip_list` attribute rather than relying solely on `private_ip`.

## Query Language Note

GES 2.5.0 and later versions no longer support Gremlin queries. Use Cypher or GQL for graph queries.
For details, see the [GES release notes](https://support.huaweicloud.com/wtsnew-ges/index.html).

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the GES graph instance may take up to 60 minutes
* `enable_https = true` may significantly impact graph query performance
* `security_group_id` is ForceNew; changing security groups requires recreating the graph instance
* When `enable_audit` is true, `audit_log_group_name` must be set to an existing LTS log group
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.50.0 |
