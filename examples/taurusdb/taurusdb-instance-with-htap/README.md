# Create a TaurusDB Instance with an HTAP StarRocks Instance

This example provides best practice code for using Terraform to create a TaurusDB instance and attach an HTAP
StarRocks instance to it in HuaweiCloud TaurusDB service.

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
* `security_group_name` - The security group name of the TaurusDB instance
* `htap_security_group_name` - The security group name of the HTAP StarRocks instance
* `taurusdb_instance_name` - The TaurusDB instance name
* `htap_instance_name` - The HTAP StarRocks instance name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `taurusdb_flavor_ref` - The flavor code of the TaurusDB instance (default: "")
* `taurusdb_root_pwd` - The database password of the TaurusDB instance (default: "")
* `taurusdb_availability_zone_mode` - The availability zone mode of the TaurusDB instance (default: "multi")
* `taurusdb_read_replicas` - The number of read replicas of the TaurusDB instance (default: 2)
* `enterprise_project_id` - The enterprise project ID (default: "0")
* `htap_db_root_pwd` - The database password of the HTAP StarRocks instance (default: "")
* `fe_flavor_id` - The specification ID of the frontend node (default: "")
* `be_flavor_id` - The specification ID of the backend node (default: "")
* `az_code` - The AZ code of the HTAP StarRocks instance (default: "")
* `fe_count` - The number of frontend nodes (default: 1)
* `be_count` - The number of backend nodes (default: 1)
* `ha_mode` - The deployment mode of the HTAP StarRocks instance (default: "Single")
* `time_zone` - The time zone of the HTAP StarRocks instance (default: "UTC+08:00")
* `enable_users_sync` - Whether to enable users synchronization (default: "true")
* `open_slow_log_switch` - Whether to enable the slow query log original text switch (default: "true")
* `volume_io_type` - The storage type of the frontend and backend nodes (default: "SSD")
* `fe_volume_capacity` - The disk capacity in GB of the frontend node (default: 50)
* `be_volume_capacity` - The disk capacity in GB of the backend node (default: 50)
* `engine_version` - The major version number of the engine (default: "")
* `be_parameter_values` - A map of parameter name and value to modify for the backend nodes
* `fe_parameter_values` - A map of parameter name and value to modify for the frontend nodes

## Usage

* Copy this example script to your `main.tf`.

* Provide the authentication variables through environment variables:

  ```bash
  $ export TF_VAR_region_name="cn-north-4"
  $ export TF_VAR_access_key="your_access_key"
  $ export TF_VAR_secret_key="your_secret_key"
  ```

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                 = "your_vpc_name"
  subnet_name              = "your_subnet_name"
  security_group_name      = "your_security_group_name"
  htap_security_group_name = "your_htap_security_group_name"
  taurusdb_instance_name   = "your_taurusdb_instance_name"
  htap_instance_name       = "your_htap_starrocks_instance_name"
  be_parameter_values      = {
    alter_tablet_worker_count            = "1"
    base_compaction_num_threads_per_disk = "1"
  }
  fe_parameter_values      = {
    alter_table_timeout_second     = "21600"
    bdbje_heartbeat_timeout_second = "10"
  }
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
* The creation of the TaurusDB instance and the HTAP StarRocks instance takes about 15-30 minutes
* This example creates the VPC, subnet, two security groups, a TaurusDB instance and the HTAP StarRocks instance
* The TaurusDB instance flavor and availability zones are automatically queried from the `huaweicloud_taurusdb_flavors`
  data source
* The HTAP StarRocks instance flavors (frontend and backend), availability zone and engine version are automatically
  queried from the `huaweicloud_taurusdb_htap_flavors` and `huaweicloud_taurusdb_htap_datastores` data sources
* The HTAP StarRocks instance password is generated randomly by default, set `htap_db_root_pwd` to use a custom one
* To deploy a cluster HTAP StarRocks instance, set `ha_mode = "Cluster"`, `fe_count = 3` and `be_count = 3`
* All resources will be created in the specified region

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.95.0 |
| random | >= 3.0.0 |
