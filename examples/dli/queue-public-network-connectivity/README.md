# Configure DLI queue connectivity with the public network

This example provides best practice code for using Terraform to configure network connectivity between a DLI
exclusive queue and the public network in HuaweiCloud via an enhanced datasource connection, NAT gateway, and SNAT.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* DLI agency `dli_management_agency` configured with **DLI Datasource Connections Agency Access**
  in the target region/project

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DLI resources are located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `elastic_resource_pool_name` - The DLI elastic resource pool name
* `elastic_resource_pool_cidr` - The CIDR of the elastic resource pool  
  This CIDR must not overlap with the VPC CIDR
* `queue_name` - The DLI exclusive queue name
* `vpc_name` - The VPC name
* `vpc_cidr` - The CIDR block of the VPC
* `subnet_name` - The subnet name
* `datasource_connection_name` - The enhanced datasource connection name
* `datasource_connection_routes` - The custom routes of the enhanced datasource connection  
  + `name` - The custom route name
  + `cidr` - The public destination CIDR to access
* `eip_bandwidth_name` - The EIP bandwidth name
* `nat_gateway_name` - The NAT gateway name

#### Optional Variables

* `elastic_resource_pool_description` - The description of the elastic resource pool (default: "")
* `elastic_resource_pool_min_cu` - The minimum CUs of the elastic resource pool (default: 16)
* `elastic_resource_pool_max_cu` - The maximum CUs of the elastic resource pool (default: 64)
* `enterprise_project_id` - The enterprise project ID (default: "")
* `queue_type` - The type of the DLI queue (default: "sql")  
  The valid values are **sql** and **general**
* `queue_cu_count` - The CU count of the DLI queue (default: 16)
* `queue_description` - The description of the DLI queue (default: "")
* `subnet_cidr` - The CIDR block of the subnet (default: "")  
  If empty, it is calculated from the VPC CIDR
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")  
  If empty, it is calculated from the subnet CIDR
* `datasource_connection_hosts` - The custom hosts of the enhanced datasource connection (default: [])  
  + `name` - The custom host name
  + `ip` - The IPv4 address of the host
* `eip_type` - The EIP type (default: "5_bgp")
* `eip_bandwidth_size` - The EIP bandwidth size in Mbps (default: 5)
* `eip_bandwidth_share_type` - The EIP bandwidth share type (default: "PER")
* `eip_bandwidth_charge_mode` - The EIP bandwidth charge mode (default: "traffic")
* `nat_gateway_spec` - The NAT gateway specification (default: "1")  
* `nat_gateway_description` - The description of the NAT gateway (default: "")
* `snat_description` - The description of the SNAT rule (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  elastic_resource_pool_name   = "your_elastic_resource_pool_name"
  elastic_resource_pool_cidr   = "your_elastic_resource_pool_cidr"
  queue_name                   = "your_queue_name"
  vpc_name                     = "your_vpc_name"
  vpc_cidr                     = "your_vpc_cidr"
  subnet_name                  = "your_subnet_name"
  datasource_connection_name   = "your_datasource_connection_name"
  datasource_connection_routes = [
    {
      name = "your_route_name"
      cidr = "your_public_destination_cidr"
    }
  ]

  eip_bandwidth_name = "your_eip_bandwidth_name"
  nat_gateway_name   = "your_nat_gateway_name"
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
* All resources will be created in the specified region
* The CIDR block of the elastic resource pool must not overlap with the VPC CIDR

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.75.5 |
