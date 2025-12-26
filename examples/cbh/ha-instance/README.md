# Create a CBH HA instance

This example provides best practice code for using Terraform to create a high-availability CBH HA instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

### Authentication Variables

* `region_name` - The region where the CBH HA instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the CBH HA instance
* `instance_password` - The login password for the CBH HA instance

#### Optional Variables

* `master_availability_zone` - The availability zone name of the master instance (default: "")
* `slave_availability_zone` - The availability zone name of the slave instance (default: "")
* `instance_flavor_id` - The flavor ID of the CBH HA instance (default: "")
* `instance_flavor_type` - The flavor type of the CBH HA instance (default: "basic")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `charging_mode` - The charging mode of the CBH HA instance (default: "prePaid")
* `period_unit` - The charging period unit of the CBH HA instance (default: "month")
* `period` - The charging period of the CBH HA instance (default: 1)
* `auto_renew` - Whether to enable auto-renew for the CBH HA instance (default: "false")

## Usage

* Copy this example script to your `main.tf`.
* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_cbh_ha_instance_name"
  instance_flavor_id  = "your_cbh_ha_instance_flavor_id"
  instance_password   = "your_cbh_ha_instance_password"
  ```

* Initialize Terraform:

   ```bash
   terraform init
   ```

* Review the Terraform plan:

   ```bash
   terraform plan
   ```

* Apply the configuration:

   ```bash
   terraform apply
   ```

* To clean up the resources:

   ```bash
   terraform destroy
   ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the CBH HA instance takes about 30 minutes
* All resources will be created in the specified region

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.64.3 |
