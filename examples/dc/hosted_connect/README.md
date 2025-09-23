# Create a DC hosted connect example

This example provides best practice code for using Terraform to create a DC (Direct Connect) hosted connect in HuaweiCloud.

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

* `hosted_connect_name` - The name of the hosted connect
* `bandwidth` - The bandwidth size of the hosted connect in Mbit/s
* `hosting_id` - The ID of the operations connection on which the hosted connect is created
* `vlan` - The VLAN allocated to the hosted connect
* `resource_tenant_id` - The tenant ID for whom a hosted connect is to be created

#### Optional Variables

* `hosted_connect_description` - The description of the hosted connect (default: "Created by Terraform")
* `peer_location` - The location of the on-premises facility at the other end of the connection (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  hosted_connect_name = "your_hosted_connect_name"
  bandwidth           = your_bandwidth
  hosting_id          = "your_hosting_id"
  vlan                = your_vlan
  resource_tenant_id  = "your_resource_tenant_id"
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

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.77.6 |
