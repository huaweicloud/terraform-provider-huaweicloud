# Create a GES graph backup

This example provides best practice code for using Terraform to create a GES (Graph Engine Service) graph backup
in HuaweiCloud. The backup depends on an existing GES graph instance, which requires VPC, subnet, and security group.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GES backup is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name for the GES graph
* `subnet_name` - The subnet name for the GES graph
* `security_group_name` - The security group name for the GES graph
* `graph_name` - The GES graph name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "192.168.0.0/24")
* `gateway_ip` - The gateway IP address of the subnet (default: "192.168.0.1")
* `graph_size_type_index` - The graph size type index (default: "1")
* `graph_cpu_arch` - The CPU architecture type of the GES graph (default: "x86_64")
* `graph_crypt_algorithm` - The cryptography algorithm of the GES graph (default: "generalCipher")
* `graph_enable_https` - Whether to enable HTTPS for the GES graph (default: false)
* `graph_tags` - The key/value pairs to associate with the GES graph
  (default: { key = "val", foo = "bar" })

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name              = "tf_test_ges_vpc"
  vpc_cidr              = "192.168.0.0/16"
  subnet_name           = "tf_test_ges_subnet"
  subnet_cidr           = "192.168.0.0/24"
  gateway_ip            = "192.168.0.1"
  security_group_name   = "tf_test_ges_secgroup"
  graph_name            = "tf_test_ges_graph"
  graph_crypt_algorithm = "generalCipher"
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
* The creation of the GES graph takes about 30 minutes, and the backup is created after the graph is ready
* The backup is automatically deleted when the associated graph is deleted
* Ensure the GES graph is in the running state before creating a backup

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.51.0 |
