# Create a CCE Autopilot Cluster

This example provides best practice code for using Terraform to create and manage a CCE autopilot cluster
in HuaweiCloud, including optional addon installation.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Cluster Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `cluster_name` - The CCE autopilot cluster name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `cluster_description` - The description of the cluster (default: "")
* `cluster_tags` - The tags of the cluster (default: {})

### Addon Variables (Optional)

* `addon_enabled` - Whether to install addon (default: false)
* `addon_name` - The name of the addon template to install (default: "coredns")
* `addon_version` - The version of the addon (default: "1.4.3")
* `addon_basic_values` - The basic values for the addon (default: {})
* `addon_flavor_values` - The flavor values for the addon (default: {})
* `addon_custom_values` - The custom values for the addon (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  cluster_name        = "your_cluster_name"
  cluster_description = "Created by terraform script"
  cluster_tags        = {
    owner = "terraform"
  }

  # Enable addon installation (optional)
  addon_enabled = true
  addon_name    = "coredns"
  addon_version = "1.4.3"
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
* The creation of the CCE autopilot cluster takes about 10-15 minutes
* This example creates the VPC, subnet, and CCE autopilot cluster
* CCE autopilot clusters use the `eni` container network mode
* Addon installation is optional and can be enabled by setting `addon_enabled = true`
* All resources will be created in the specified region

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.90.0 |
