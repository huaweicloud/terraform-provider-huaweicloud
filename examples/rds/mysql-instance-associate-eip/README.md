# Create a MySQL RDS instance with EIP association

This example provides best practice code for using Terraform to create a MySQL RDS instance in HuaweiCloud RDS service
and associate it with an Elastic IP (EIP) for external access.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the MySQL RDS instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The MySQL RDS instance name
* `instance_backup_time_window` - The backup time window in HH:MM-HH:MM format
* `instance_backup_keep_days` - The number of days to retain backups
* `bandwidth_name` - The name for the bandwidth (required if EIP address is not provided)

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `availability_zone` - The availability zone to which the RDS instance belongs (default: "")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `instance_flavor_id` - The flavor ID of the RDS instance (default: "")
* `instance_db_type` - The database engine type (default: "MySQL")
* `instance_db_version` - The database engine version (default: "8.0")
* `instance_db_port` - The database port (default: 3306)
* `instance_password` - The password for the RDS instance (default: "")
* `instance_mode` - The instance mode for the RDS instance flavor (default: "single")
* `instance_flavor_group_type` - The group type for the RDS instance flavor (default: "general")
* `instance_flavor_vcpus` - The number of the RDS instance CPU cores for the RDS instance flavor (default: 2)
* `instance_volume_type` - The storage volume type (default: "CLOUDSSD")
* `instance_volume_size` - The storage volume size in GB (default: 40)
* `associate_eip_address` - The EIP address to associate with the RDS instance (default: "")
* `eip_type` - The type of the EIP (default: "5_bgp")
* `bandwidth_size` - The size of the bandwidth (default: 5)
* `bandwidth_share_type` - The share type of the bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")

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
  bandwidth_name              = "your_bandwidth_name"
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
2. **EIP Association**: Associates an Elastic IP with the RDS instance for external access
3. **Flexible EIP Configuration**: Supports both creating new EIP and using existing EIP
4. **Network Configuration**: Sets up VPC, subnet, and security group for the RDS instance
5. **Backup Strategy**: Configures automated backup with customizable time window and retention period

## EIP Association Options

### Option 1: Create New EIP

If you don't provide an existing EIP address, the example will create a new EIP with the specified bandwidth configuration:

```hcl
bandwidth_name = "your_bandwidth_name"
bandwidth_size = 10
eip_type       = "5_bgp"
```

### Option 2: Use Existing EIP

If you have an existing EIP, you can associate it directly:

```hcl
associate_eip_address = "your_existing_eip_address"
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the MySQL RDS instance takes about 5 minutes
* This example creates the MySQL RDS instance, VPC, subnet, security group, EIP, and EIP association
* The EIP association enables external access to the RDS instance
* All resources will be created in the specified region
* Ensure your security group rules allow access from the external network if needed

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.54.0 |
| random | >= 3.0.0 |
