# Create a GaussDB disaster recovery architecture

This example provides best practice code for using Terraform to create a GaussDB disaster recovery
architecture in HuaweiCloud. It deploys a primary instance and a standby instance in different
availability zones within the same region, configures the DR settings, and establishes the
DR relationship.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GaussDB instances are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Required Variables

* `primary_vpc_name` - The VPC name in the primary region
* `primary_subnet_names` - The subnet names for each AZ in the primary region
* `primary_availability_zones` - The availability zones in the primary region
* `primary_security_group_name` - The security group name in the primary region
* `dr_vpc_name` - The VPC name in the DR region
* `dr_subnet_name` - The subnet name in the DR region
* `dr_availability_zone` - The availability zone in the DR region
* `dr_security_group_name` - The security group name in the DR region
* `primary_instance_name` - The name of the primary GaussDB instance
* `instance_flavor` - The spec_code of the GaussDB instance flavor
* `primary_instance_availability_zones` - The comma-separated AZ string for the primary GaussDB instance
* `enterprise_project_id` - The enterprise project ID to which the GaussDB instances belong
* `dr_instance_name` - The name of the DR GaussDB instance
* `dr_instance_availability_zones` - The comma-separated AZ string for the DR GaussDB instance
* `dr_user_name` - The database username for disaster recovery
* `dr_user_password` - The database password for disaster recovery

### Optional Variables

* `primary_vpc_cidr` - The CIDR block of the VPC in the primary region (default: "172.16.0.0/16")
* `secgroup_rules` - The security group ingress rules (default: rules for DR ports)
* `dr_vpc_cidr` - The CIDR block of the VPC in the DR region (default: "172.17.0.0/16")
* `dr_region_name` - The standby region where the DR GaussDB instance is located (default: "cn-east-3")
* `cc_bandwidth` - The inter-region bandwidth (Mbit/s) for the Cloud Connection (default: 10)
* `instance_passwords` - The passwords for GaussDB instances, empty string to auto-generate (default: ["", ""])
* `instance_db_port` - The database port (default: 5432)
* `primary_instance_volume_type` - The storage volume type of the primary GaussDB instance (default: "ULTRAHIGH")
* `primary_instance_volume_size` - The storage volume size (GB) of the primary GaussDB instance (default: 40)
* `dr_instance_volume_type` - The storage volume type of the DR GaussDB instance (default: "ULTRAHIGH")
* `dr_instance_volume_size` - The storage volume size (GB) of the DR GaussDB instance (default: 40)
* `dr_disaster_type` - The disaster recovery type (default: "stream")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  primary_vpc_name                    = "example_vpc"
  primary_vpc_cidr                    = "172.16.0.0/16"
  primary_subnet_names                = ["example_subnet_1", "example_subnet_2"]
  primary_availability_zones          = ["primary_region_az_1", "primary_region_az_2"]
  primary_security_group_name         = "example_security_group"
  dr_vpc_name                         = "example_vpc_dr"
  dr_vpc_cidr                         = "172.17.0.0/16"
  dr_subnet_name                      = "example_subnet_dr"
  dr_availability_zone                = "dr_region_az"
  dr_security_group_name              = "example_security_group"
  dr_region_name                      = "dr_region"
  instance_passwords                  = ["password_1", "password_2"]
  primary_instance_name               = "example_primary_instance"
  instance_flavor                     = "gaussdb.opengauss.ee.c3.xlarge.x864.ha"
  primary_instance_availability_zones = "primary_region_az_1,primary_region_az_2,primary_region_az_3"
  instance_db_port                    = 5432
  enterprise_project_id               = "your_enterprise_project_id"
  primary_instance_volume_type        = "ULTRAHIGH"
  primary_instance_volume_size        = 40
  dr_instance_name                    = "example_dr_instance"
  dr_instance_availability_zones      = "dr_region_az_1,dr_region_az_2,dr_region_az_3"
  dr_instance_volume_type             = "ULTRAHIGH"
  dr_instance_volume_size             = 40
  dr_user_name                        = "root"
  dr_user_password                    = "your_dr_user_password"
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

1. **GaussDB DR Architecture**: Deploys primary and standby GaussDB instances in different
   availability zones within the same region
2. **DR Configuration Reset**: Configures the opposite data CIDR for both instances before
   establishing the DR relationship
3. **DR Relationship**: Establishes the disaster recovery relationship between the primary
   and standby instances
4. **Network Configuration**: Sets up VPC, subnets, and security group for the instances
5. **Centralized Mode**: Uses centralized HA mode with strong consistency and synchronous
   replication

## DR Operations Reference

After the DR architecture is deployed, you can perform the following operations using
corresponding Terraform resources:

| Operation | Resource | Description |
|-----------|----------|-------------|
| Failover | `huaweicloud_gaussdb_dr_instance_to_primary` | Promote standby to primary in emergency |
| Switchover | `huaweicloud_gaussdb_dr_instance_primary_role_switch` | Planned primary/standby role switch |
| DR Drill | `huaweicloud_gaussdb_dr_drill` | Simulate disaster recovery drill |
| Log Cache | `huaweicloud_gaussdb_dr_log_cache` | Keep DR logs for recovery |
| Reestablish | `huaweicloud_gaussdb_dr_relationship_reestablish` | Rebuild broken DR relationship |
| Monitor | `data.huaweicloud_gaussdb_instance_dr_status` | Query RPO/RTO status |

## Cross-Region DR

This example deploys DR architecture within the same region. To deploy cross-region DR,
add a second provider with a different region and update the standby instance to use it:

```hcl
provider "huaweicloud" {
  alias      = "standby"
  region     = var.dr_region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
```

Then create the standby VPC, subnet, and instance under the `huaweicloud.standby` provider.

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of GaussDB instances takes about 10-15 minutes each
* Both primary and standby instances must have the same flavor and configuration
* The DR configuration reset must be performed before establishing the DR relationship
* All resources will be created in the specified region

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0  |
| huaweicloud | >= 1.66.3 |
