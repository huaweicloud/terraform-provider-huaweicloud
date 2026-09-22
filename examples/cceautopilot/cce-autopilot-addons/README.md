# Deploy a CCE Autopilot Addon

This example provides best practice code for using Terraform to create a CCE autopilot cluster
and deploy an addon (log-agent) in HuaweiCloud.

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

### VPC and Subnet Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")

### Cluster Variables

#### Required Variables

* `cluster_name` - The CCE autopilot cluster name

#### Optional Variables

* `cluster_version` - The Kubernetes version of the cluster (default: "v1.36")
* `cluster_alias` - The alias of the cluster (default: "")
* `cluster_description` - The description of the cluster (default: "CCE Autopilot cluster for addon deployment")
* `cluster_category` - The cluster type (default: "Turbo")
* `cluster_type` - The master node architecture (default: "VirtualMachine")
* `cluster_custom_san` - The custom SAN field in the API server certificate (default: [])
* `cluster_eip_id` - The EIP ID of the cluster (default: "")
* `cluster_enable_snat` - Whether SNAT is configured for the cluster (default: false)
* `cluster_enable_swr_image_access` - Whether SWR image access is enabled (default: true)
* `cluster_delete_efs` - Whether to delete the associated EFS when deleting the cluster (default: false)
* `cluster_delete_eni` - Whether to delete the associated ENI when deleting the cluster (default: false)
* `cluster_delete_net` - Whether to delete the associated network resources (default: false)
* `cluster_delete_obs` - Whether to delete the associated OBS when deleting the cluster (default: false)
* `cluster_delete_sfs_turbo` - Whether to delete the associated SFS Turbo (default: false)
* `cluster_lts_reclaim_policy` - The LTS reclaim policy (default: "retain")
* `cluster_enterprise_project_id` - The ID of the enterprise project (default: "0")
* `cluster_tags` - The tags of the cluster (default: {})

### SWR Organization Variables

#### Required Variables

* `swr_organization_name` - The name of the SWR organization

### Addon Variables

#### Optional Variables

* `addon_template_name` - The name of the add-on template (default: "log-agent")
* `addon_version` - The version of the add-on. If not specified, the latest compatible version will be used (default: "")
* `addon_name` - The name of the add-on (default: "")
* `addon_alias` - The alias of the add-on (default: "")
* `addon_values_basic` - The basic configuration values for the add-on (default: {})
* `addon_values_flavor` - The flavor configuration values for the add-on (default: {})
* `addon_values_custom` - The custom configuration values for the add-on (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  region_name           = "cn-north-4"
  vpc_name              = "your_vpc_name"
  subnet_name           = "your_subnet_name"
  cluster_name          = "your_cluster_name"
  swr_organization_name = "your_swr_organization_name"
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
* This example creates a VPC, subnet, CCE autopilot cluster, SWR organization, and deploys a log-agent addon
* The creation of the CCE autopilot cluster takes about 10-15 minutes
* CCE autopilot clusters use the `eni` container network mode
* The SWR organization name must be globally unique across HuaweiCloud
* The log-agent addon requires SWR image access to be enabled on the cluster
* All resources will be created in the specified region
* The addon version is automatically retrieved from the addon template data source

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.72.0 |
| random | >= 3.0.0 |
