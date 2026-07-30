# Create a GaussDB data backup

This example provides best practice code for using Terraform to create a GaussDB data backup
in HuaweiCloud. It deploys a GaussDB instance and creates a manual backup for it.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GaussDB instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The name of the GaussDB instance
* `backup_name` - The name for the manual backup

### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: `172.16.0.0/16`)
* `subnet_cidr` - The CIDR block of the subnet. If empty, it will be calculated from the VPC CIDR (default: `""`)
* `subnet_gateway_ip` - The gateway IP of the subnet.
  If empty, it will be calculated from the subnet CIDR (default: `""`)
* `security_group_rule_ports` - The security group ingress rule ports
  (default: `2379-2380,5000-5001,5432-5532,6000,6500,12016,20050`)
* `instance_password` - The password for the GaussDB instance.
  If empty, a random password will be generated (sensitive, default: `""`)
* `instance_flavor` - The flavor of the GaussDB instance (default: `gaussdb.opengauss.ee.c3.xlarge.x864.ha`)
* `instance_availability_zones` - The availability zones for the GaussDB instance, comma-separated.
  If not specified, the first 3 available zones will be used automatically (default: `""`)
* `instance_db_port` - The database port of the GaussDB instance (default: `5432`)
* `enterprise_project_id` - The enterprise project ID of the GaussDB instance (default: `null`)
* `instance_ha_mode` - The HA mode of the GaussDB instance (default: `centralization_standard`)
* `instance_ha_replication_mode` - The HA replication mode of the GaussDB instance (default: `sync`)
* `instance_ha_consistency` - The HA consistency of the GaussDB instance (default: `strong`)
* `instance_volume_type` - The storage volume type of the GaussDB instance (default: `ULTRAHIGH`)
* `instance_volume_size` - The storage volume size (GB) of the GaussDB instance (default: `40`)
* `backup_description` - The description for the manual backup (default: `""`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_gaussdb_instance_name"
  backup_name         = "your_manual_backup_name"
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

## Backup Operations Reference

After the backup is created, you can perform the following operations using
corresponding Terraform resources:

| Operation | Resource | Description |
|-----------|----------|-------------|
| Create Backup | `huaweicloud_gaussdb_backup` | Create a manual backup for the instance |
| Stop Backup | `huaweicloud_gaussdb_backup_stop` | Stop an in-progress backup task |
| Restore from Backup | `huaweicloud_gaussdb_restore` | Restore instance from a backup or point in time |
| List Backups | `data.huaweicloud_gaussdb_backups` | Query the backup list of an instance |
| Backup Config | `data.huaweicloud_gaussdb_backup_configurations` | Query backup configuration details |
| Backup Files | `data.huaweicloud_gaussdb_backup_files` | Query backup file download links |

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the GaussDB instance takes about 10-15 minutes
* The backup creation time depends on the data volume
* The instance password defaults to an empty string, which will auto-generate a random password
  using the `random_password` resource. You can also set a custom password via `instance_password`
* The GaussDB instance uses centralization_standard HA mode with sync replication and strong consistency by default,
  configurable via `instance_ha_mode`, `instance_ha_replication_mode`, and `instance_ha_consistency`
* The `lifecycle.ignore_changes` for `flavor` prevents unintended instance specification changes during updates
* All resources will be created in the specified region

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.90.0 |
