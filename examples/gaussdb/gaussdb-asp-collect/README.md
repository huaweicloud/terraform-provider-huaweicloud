# GaussDB ASP collect

This example provides best practice code for using Terraform to collect ASP (Application Support Package) data from a
GaussDB instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing GaussDB instance

## Required Variables

### Authentication Variables

* `region_name` - The region where the GaussDB ASP collect resource is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The GaussDB instance name
* `start_time` - The start time for ASP collect in `yyyy-mm-ddThh:mm:ssZ` format
* `end_time` - The end time for ASP collect in `yyyy-mm-ddThh:mm:ssZ` format

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
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_gaussdb_instance_name"
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

* The time format is `yyyy-mm-ddThh:mm:ssZ`, where `Z` indicates the timezone offset (e.g., `+0800`)
* The resource uses the GaussDB (opengauss) service endpoint
* The create operation sends a POST request and waits for the job to complete (default timeout: 30 minutes)
* The resource ID is the job ID returned from the create API response
* All input parameters (`instance_id`, `start_time`, `end_time`) are non-updatable and will trigger resource replacement
  if changed
* The update operation is a no-op
* Deleting this resource only removes it from Terraform state, the ASP report remains in the cloud
* The resource supports import using `<instance_id>/<id>` format
* The `download_url` is valid for 30 minutes after collection
* The computed attributes include `file_size`, `download_url`, `status`, `file_path`, `file_name`, and `obs_bucket`
* The `status` field can be `SUCCESS`, `FAILED`, or `EXPORTING`

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.92.0 |
