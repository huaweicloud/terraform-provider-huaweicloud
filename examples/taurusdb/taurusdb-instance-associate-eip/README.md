# Create a TaurusDB Instance and Associate EIP

This example provides best practice code for using Terraform to create a TaurusDB instance and associate an EIP for
public network access in HuaweiCloud TaurusDB service.

It supports both creating a new EIP and using an existing EIP by configuring the `associate_eip_address` variable.

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

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The TaurusDB instance name
* `instance_backup_time_window` - The backup time window in HH:MM-HH:MM format
* `instance_backup_keep_days` - The number of days to retain backups

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `availability_zone_mode` - The availability zone mode (default: "multi")
* `master_availability_zone` - The master availability zone (default: "")
* `instance_db_port` - The database port (default: 3306)
* `instance_password` - The password for the TaurusDB instance (default: "")
* `instance_flavor_ref` - The flavor code of the TaurusDB instance (default: "")
* `instance_mode` - The instance mode (default: "Cluster")
* `read_replicas` - The number of read replicas (default: 2)
* `enterprise_project_id` - The enterprise project ID (default: "0")
* `associate_eip_address` - The existing EIP address to associate.
  If not specified, a new EIP will be created (default: "")
* `eip_type` - The EIP type (default: "5_bgp")
* `bandwidth_name` - The bandwidth name (required when creating new EIP)
* `bandwidth_size` - The bandwidth size in Mbit/s (default: 5)
* `bandwidth_share_type` - The share type of the bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  **Create a new EIP:**

  ```hcl
  vpc_name                    = "your_vpc_name"
  subnet_name                 = "your_subnet_name"
  security_group_name         = "your_security_group_name"
  instance_name               = "your_taurusdb_instance_name"
  instance_backup_time_window = "02:00-03:00"
  instance_backup_keep_days   = 7
  bandwidth_name              = "your_bandwidth_name"
  ```

  **Use an existing EIP:**

  ```hcl
  vpc_name               = "your_vpc_name"
  subnet_name            = "your_subnet_name"
  security_group_name    = "your_security_group_name"
  instance_name          = "your_taurusdb_instance_name"
  instance_backup_time_window = "02:00-03:00"
  instance_backup_keep_days   = 7
  associate_eip_address  = "your_existing_eip_address"
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
* The creation of the TaurusDB instance takes about 15-20 minutes
* When `associate_eip_address` is not specified, a new EIP will be created and associated to the TaurusDB instance
* When `associate_eip_address` is specified, the existing EIP will be queried and associated to the TaurusDB instance
* The `bandwidth_name` is required when creating a new EIP (i.e., when `associate_eip_address` is not specified)
* All resources will be created in the specified region

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 1.9.0  |
| huaweicloud | >= 1.91.0 |
| random | >= 3.0.0  |
