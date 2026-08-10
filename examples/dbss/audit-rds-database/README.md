# Create a DBSS audit RDS database

This example provides best practice code for using Terraform to create a DBSS instance to audit RDS database
  in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DBSS instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `rds_instance_name` - The RDS instance name
* `dbss_instance_name` - The DBSS instance name

#### Optional Variables

* `availability_zone` - The availability zone to which the DBSS instance belongs (default: "")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `rds_instance_flavor` - The flavor of the RDS instance (default: "")
* `database_type` - The database type of the RDS instance (default: "MySQL")
* `database_version` - The database version of the RDS instance (default: "8.0")
* `instance_mode` - The mode of the RDS instance (default: "single")
* `instance_group_type` - The performance specification (default: "dedicated")
* `instance_flavor_vcpus` - The number of vCPUs of the RDS instance flavor (default: 4)
* `volume_type` - The type of the volume (default: "CLOUDSSD")
* `volume_size` - The size of the volume in GB (default: 100)
* `dbss_instance_flavor` - The flavor ID of the DBSS instance (default: "")
* `instance_spec_code` - The spec code of the DBSS instance (default: "dbss.bypassaudit.low")
* `instance_description` - The description of the DBSS instance (default: "")
* `instance_tags` - The tags of the DBSS instance (default: {})
* `enterprise_project_id` - The enterprise project ID  (default: null)
* `charging_mode` - The charging mode of the DBSS instance (default: "prePaid")
* `period_unit` - The period unit of the DBSS instance, only required when `charging_mode` is `prePaid` (default: "month")
* `period` - The period of the DBSS instance, only required when `charging_mode` is `prePaid` (default: 1)
* `auto_renew` - The auto renew of the DBSS instance (default: "false")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  rds_instance_name   = "your_rds_instance_name"
  dbss_instance_name  = "your_dbss_instance_name"
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
* The creation of the DBSS instance takes about 20 minutes
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.71.0 |
