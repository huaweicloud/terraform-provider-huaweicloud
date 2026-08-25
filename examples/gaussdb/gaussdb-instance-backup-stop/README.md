# GaussDB backup stop

This example provides best practice code for using Terraform to stop an ongoing backup of a GaussDB OpenGauss instance
in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing GaussDB OpenGauss instance with an ongoing backup

## Required Variables

### Authentication Variables

* `region_name` - The region where the GaussDB backup stop resource is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The GaussDB instance name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: `"192.168.0.0/16"`)
* `enterprise_project_id` - The ID of the enterprise project (default: `""`)
* `subnet_cidr` - The CIDR block of the subnet (default: `""`, auto-calculated from VPC CIDR)
* `gateway_ip` - The gateway IP address of the subnet (default: `""`, auto-calculated)
* `instance_password` - The password for the GaussDB instance (default: `""`, auto-generated)
* `instance_volume_type` - The storage volume type (default: `"ULTRAHIGH"`)
* `instance_volume_size` - The storage volume size in GB (default: `40`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name               = "your_vpc_name"
  subnet_name            = "your_subnet_name"
  security_group_name    = "your_security_group_name"
  instance_name          = "your_gaussdb_instance_name"
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

## Note

* This resource is an action resource that stops an ongoing backup of a GaussDB OpenGauss instance
* The resource uses the GaussDB (opengauss) service endpoint
* The create operation sends a POST request to `/v3/{project_id}/instances/{instance_id}/backups/stop` and waits for the
  job to complete (default timeout: 30 minutes)
* The resource ID equals the `instance_id`
* The `instance_id` parameter is non-updatable and will trigger resource replacement if changed (ForceNew)
* The read operation is a no-op (returns nil), no API call is made during read
* There is no update operation
* Deleting this resource only removes it from Terraform state with a warning, the GaussDB OpenGauss instance remains in
  the cloud
* This resource does not support import

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.90.0 |
