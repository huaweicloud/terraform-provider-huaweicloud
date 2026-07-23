# Create a TaurusDB Instance with Database Proxy

This example provides best practice code for using Terraform to create a TaurusDB instance with a database proxy in
HuaweiCloud TaurusDB service.

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

#### Optional Variables - TaurusDB Instance

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `master_availability_zone` - The master availability zone (default: "")
* `instance_flavor_ref` - The flavor code (default: "")
* `instance_password` - The password for the TaurusDB instance (default: "")
* `read_replicas` - The number of read replicas (default: 4)
* `enterprise_project_id` - The enterprise project ID (default: "0")
* `instance_db_port` - The database port (default: 3306)
* `volume_type` - The storage type of the instance (default: "DL6")
* `time_zone` - The time zone of the instance (default: "UTC+08:00")
* `ssl_option` - Whether to enable SSL (default: "true")
* `sql_filter_enabled` - Whether to enable SQL filter (default: true)
* `slow_log_show_original_switch` - Whether to enable slow log show original switch (default: true)
* `table_name_case_sensitivity` - Whether the kernel table name is case sensitive (default: true)
* `multi_tenant_switch` - Whether to enable multi-tenancy switch (default: "true")
* `maintain_begin` - The start time of the maintenance window (default: "08:00")
* `maintain_end` - The end time of the maintenance window (default: "11:00")
* `description` - The description of the TaurusDB instance (default: "")
* `seconds_level_monitoring_enabled` - Whether to enable seconds level monitoring (default: true)
* `seconds_level_monitoring_period` - The seconds level collection period (default: 5)
* `audit_log_enabled` - Whether to enable audit log (default: true)
* `audit_log_keep_days` - The number of days for storing audit logs (default: 7)
* `reserve_audit_logs` - Whether to reserve historical audit logs when SQL audit is disabled (default: "true")
* `tags` - The tags of the TaurusDB instance (default: {})

#### Optional Variables - Database Proxy

* `proxy_mode` - The type of the proxy (default: "readwrite")
* `proxy_node_num` - The number of proxy nodes (default: 2)
* `proxy_name` - The name of the database proxy (default: "")
* `route_mode` - The routing policy of the proxy (default: 1)
* `proxy_new_node_auto_add_status` - Whether to automatically add new nodes (default: "OFF")
* `proxy_new_node_weight` - The weight of new nodes (default: 20)
* `proxy_port` - The proxy port (default: 3306)
* `proxy_transaction_split` - Whether to enable transaction split (default: "OFF")
* `proxy_consistence_mode` - The consistency mode (default: "eventual")
* `proxy_connection_pool_type` - The connection pool type (default: "CLOSED")
* `proxy_open_access_control` - Whether to enable access control (default: true)
* `access_control_type` - The access control mode (default: "white")
* `proxy_dns_name_prefix` - The DNS name prefix for the proxy (default: "")
* `proxy_master_node_weight` - The weight of the master node (default: 20)
* `proxy_readonly_node_weight` - The weight of read-only nodes (default: 30)
* `proxy_access_control_ip_list` - The access control IP list (default: [])
* `proxy_parameters` - The parameters for the proxy (default: [])

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                       = "your_vpc_name"
  subnet_name                    = "your_subnet_name"
  security_group_name            = "your_security_group_name"
  instance_name                  = "your_taurusdb_instance_name"
  instance_backup_time_window    = "02:00-03:00"
  instance_backup_keep_days      = 7
  proxy_name                     = "your_proxy_name"
  proxy_node_num                 = 2
  proxy_port                     = 3339
  proxy_access_control_ip_list   = [
    {
      ip          = "3.3.3.3"
      description = "test description"
    }
  ]
  proxy_parameters               = [
    {
      name      = "multiStatementType"
      value     = "Loose"
      elem_type = "system"
    }
  ]
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
* The creation of the TaurusDB instance takes about 20-30 minutes
* The TaurusDB instance must be deployed in multi-AZ mode to use the database proxy
* It is recommended to have at least 4 read replicas when using the database proxy
* Node weights are assigned automatically based on the instance nodes list
* The instance flavor and availability zones are automatically queried from `huaweicloud_taurusdb_flavors` data source
* The proxy flavor is automatically queried from `huaweicloud_taurusdb_proxy_flavors` data source

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.95.0|
| random | >= 3.0.0 |
