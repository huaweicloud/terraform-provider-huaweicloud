# MySQL instance with public IP

This example provides best practice code for using Terraform to create a configurable MySQL RDS instance on
HuaweiCloud with VPC networking, security group, and eip binding.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region name
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `secgroup_name` - The name of security group
* `instance_name` - The name of the RDS instance
* `backup_time_window` - The backup time window in HH:MM-HH:MM format
* `backup_keep_days` - The number of days to retain backups

#### Optional Variables

* `availability_zone` - The availability zone (default: "")
* `flavor_id` - The flavor ID for the instance (default: "")
* `db_type` - The database engine type (default: "MySQL")
* `db_version` - The database engine version (default: "8.0")
* `db_port` -  The database port (default: 3306)
* `instance_mode` - The instance mode for the RDS instance flavor (default: "single")
* `group_type` - The group type (default: "general")
* `vcpus` - The CPU of flavor for the instance (default: 2)
* `vpc_id` - The ID of the existing VPC (default: "")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_id` - The ID of the existing subnet (default: "")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway` - The gateway IP of the subnet (default: "")
* `secgroup_id` - The ID of the existing security group (default: "")
* `protocol` - The protocol type for network security rule (default: "tcp")
* `charging_mode` - The billing method (default: "postPaid")
* `volume_type` - The storage volume type (default: "CLOUDSSD")
* `volume_size` - The storage volume size in GB (default: 40)
* `primary_dns` - Private DNS server IP address (default: "100.125.1.250")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name           = "your_vpc_name"
  subnet_name        = "your_subnet_name"
  secgroup_name      = "your_security_group_name"
  instance_name      = "your_instance_name"
  backup_time_window = "08:00-09:00"
  backup_keep_days   = 1
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

* Passwords are auto-generated with special characters and stored in Terraform state
* Please read the implicit and explicit dependencies in the script carefully
* All resources will be created in the specified region

## Requirements

| Name        | Version    |
|-------------|------------|
| terraform   | >= 0.12.0  |
| huaweicloud | >= 1.67.0  |
| random      | >= 3.7.2   |
