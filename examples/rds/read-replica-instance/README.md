# Create a MySQL RDS instance with read replica

This example provides best practice code for using Terraform to create a MySQL RDS instance with read replica in
HuaweiCloud RDS service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the PostgreSQL RDS instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The name of the RDS instance
* `instance_backup_time_window` - The backup time window in HH:MM-HH:MM format
* `instance_backup_keep_days` - The number of days to retain backups
* `replica_instance_name` - The name of the read replica instance

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `availability_zones` - The list of availability zones to which the RDS instance belong (default: [])
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `instance_flavors_filter` - The filter configuration of the RDS MySQL flavor instance mode for the RDS instance
  (default: [{"instance_mode": "ha"}, {"instance_mode": "replica"}])
* `instance_db_port` - The database port (default: 5432)
* `instance_password` - The password for the RDS instance (default: "")
* `instance_flavor_id` - The flavor ID of the RDS instance (default: "")
* `ha_replication_mode` - The HA replication mode of the RDS instance (default: "async")
* `instance_volume_type` - The storage volume type (default: "CLOUDSSD")
* `instance_volume_size` - The storage volume size in GB (default: 40)
* `replica_instance_flavor_id` - The flavor ID of the read replica instance (default: "")
* `replica_instance_volume_type` - The storage volume type (default: "CLOUDSSD")
* `replica_instance_volume_size` - The storage volume size in GB (default: 40)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                    = "your_vpc_name"
  subnet_name                 = "your_subnet_name"
  security_group_name         = "your_security_group_name"
  instance_name               = "your_mysql_instance_name"
  instance_backup_time_window = "08:00-09:00"
  instance_backup_keep_days   = 1
  replica_instance_name       = "your_replica_instance_name"
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

## Features

This example demonstrates the following features:

1. **MySQL RDS Instance Creation**: Creates a complete MySQL RDS instance with all necessary components
2. **Read Replica Instance**: Creates a read replica instance for the primary RDS instance
3. **Network Configuration**: Sets up VPC, subnet, and security group for the RDS instance
4. **Backup Strategy**: Configures automated backup with customizable time window and retention period
5. **High Availability**: Configures HA replication mode for the primary instance

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the MySQL RDS instance and read replica takes about 10 minutes
* This example creates the MySQL RDS instance, read replica instance, VPC, subnet, and security group
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.40.0 |
| random | >= 3.0.0 |
