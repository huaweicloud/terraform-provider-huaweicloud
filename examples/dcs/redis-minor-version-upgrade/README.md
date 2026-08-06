# DCS Instance Minor Version Upgrade

This example provides best practice code for using Terraform to upgrade the minor version of a
DCS (Distributed Cache Service) instance in HuaweiCloud. It creates a Redis HA instance, then
triggers a minor version upgrade for both the engine and proxy components.

## Prerequisites

* A HuaweiCloud account with DCS permissions
* An existing VPC and subnet (queried by name via data sources, or specify IDs directly)
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DCS instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `instance_name` - The name of the Redis single instance

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `availability_zone` - The availability zone to which the Redis single instance belongs (default: "")
* `instance_flavor_id` - The flavor ID of the Redis single instance (default: "")
* `instance_capacity` - The capacity of the Redis instance (default: 1)
* `instance_engine_version` - The engine version of the Redis single instance (default: "7.0")
* `instance_password` - The password for the Redis instance (default: null)
* `engine_minor_version` - The target engine minor version. Use `latest` for the latest version (default: `"latest"`)
* `proxy_minor_version` - The target proxy minor version. Use `latest` for the latest version (default: `"latest"`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name      = "your_vpc_name"
  subnet_name   = "your_subnet_name"
  instance_name = "your_instance_name"
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

### Action Resource

This resource is an **action** resource. Its behavior differs from standard Terraform resources:

* **Create triggers the upgrade**: The `create` operation sends a `POST` request to the DCS API
  (`POST /v2/{project_id}/instances/{instance_id}/minor-version/upgrade`) to upgrade the minor
  version of the DCS instance. The API returns a `task_id` on success. The resource is tracked in
  Terraform state with a randomly generated UUID as its ID.
* **Read is a no-op**: The `read` operation does not query the API. Terraform will not detect
  configuration drift for this resource.
* **Update is a no-op**: The `update` operation does nothing. All three parameters (`instance_id`,
  `engine_minor_version`, `proxy_minor_version`) are non-updatable — changing any of them will
  trigger resource replacement (ForceNew).
* **Delete is a no-op**: Running `terraform destroy` only removes the resource from Terraform state.
  The upgrade is not rolled back. A warning message is displayed during destruction.
* **No import support**: This resource cannot be imported into Terraform state.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.92.0 |
