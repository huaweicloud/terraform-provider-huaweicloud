# Manage GaussDB(DWS) workload resources

This example provides best practice code for using Terraform to manage GaussDB(DWS) workload
resources in HuaweiCloud, including a cluster, workload queue, cluster user (with attributes
and optional grants), queue user association, workload plan with stage, and exception rule.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DWS resources are located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `vpc_cidr` - The CIDR block of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `cluster_name` - The name of the DWS cluster
* `cluster_admin_user_name` - The administrator username of the DWS cluster
* `cluster_admin_user_pwd` - The administrator password of the DWS cluster
* `workload_queue_name` - The name of the workload queue
* `workload_queue_configurations` - The configurations of the workload queue
  + `resource_name` - The resource name
  + `resource_value` - The resource value
* `user_name` - The name of the cluster user
* `user_password` - The password of the cluster user
* `workload_plan_name` - The name of the workload plan
* `workload_plan_stage_name` - The name of the workload plan stage
* `workload_plan_stage_configurations` - The configurations of the workload plan stage  
  + `resource_name` - The resource name
  + `resource_value` - The resource value
  + `value_unit` - The unit of the resource value (optional)
  + `resource_description` - The description of the resource (optional)
* `exception_rule_name` - The name of the cluster exception rule
* `exception_rule_configurations` - The configurations of the exception rule
  + `key` - The configuration key
  + `value` - The configuration value

#### Optional Variables

* `availability_zone` - The availability zone of the DWS cluster (default: "")
* `enterprise_project_id` - The ID of the enterprise project (default: "")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `security_group_delete_default_rules` - Whether to delete the default rules of the security group (default: true)
* `security_group_rule_ports` - The security group ingress rule ports for DWS (default: "8000-10000")
* `cluster_node_type` - The flavor of the DWS cluster node (default: "")
* `cluster_version` - The version of the DWS cluster (default: "")
* `cluster_vcpus` - The vcpus of the DWS cluster (default: 4)
* `cluster_memory` - The memory of the DWS cluster (default: 32)
* `cluster_datastore_type` - The datastore type of the DWS cluster (default: "dws")
* `cluster_number_of_node` - The number of nodes in the DWS cluster (default: 3)
* `cluster_number_of_cn` - The number of CN nodes in the DWS cluster (default: 3)
* `cluster_volume_type` - The volume type of the DWS cluster (default: "SSD")
* `cluster_volume_capacity` - The volume capacity of the DWS cluster in GB (default: "100")
* `user_description` - The description of the cluster user (default: "")
* `user_cascade` - Whether to cascade delete dependencies when deleting the user or role (default: true)
* `user_login` - Whether to allow the user to log in (default: true)
* `user_create_role` - Whether to grant the permission to create roles (default: true)
* `user_create_db` - Whether to grant the permission to create databases (default: true)
* `user_system_admin` - Whether to grant the system administrator permission (default: null)
* `user_audit_admin` - Whether to grant the audit administrator permission (default: null)
* `user_inherit` - Whether to inherit permissions from roles (default: true)
* `user_use_ft` - Whether to grant the external table permission (default: null)
* `user_conn_limit` - The maximum number of concurrent connections. -1 means unlimited (default: -1)
* `user_replication` - Whether to grant the replication permission (default: null)
* `user_valid_begin` - The valid begin time of the cluster user (default: "")
* `user_valid_until` - The valid until time of the cluster user (default: "")
* `user_grant_list` - The set of grants for the cluster user (default: [])  
  + `type` - The object type
  + `database` - The database name (optional)
  + `schema_name` - The schema name (optional)
  + `object_name` - The object name (optional)
  + `all_object` - Whether all objects are effective (optional)
  + `future` - Whether future objects are effective (optional)
  + `future_object_owners` - The owners of future objects (optional)
  + `column_names` - The set of column names (optional)  
  + `privileges` - The set of privileges
    - `permission` - The privilege name
    - `grant_with` - Whether the grant option is included
* `workload_plan_stage_month` - The month of the workload plan stage (default: null)  
  If not specified, it will be applied to all months
* `workload_plan_stage_day` - The day of the workload plan stage (default: null)  
  If not specified, it will be applied to all days
* `workload_plan_stage_start_time` - The start time of the workload plan stage. The format is hh:mm:ss (default: "00:00:00")
* `workload_plan_stage_end_time` - The end time of the workload plan stage. The format is hh:mm:ss (default: "23:59:59")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                      = "your_vpc_name"
  vpc_cidr                      = "your_vpc_cidr"
  subnet_name                   = "your_subnet_name"
  security_group_name           = "your_security_group_name"
  cluster_name                  = "your_dws_cluster_name"
  cluster_admin_user_name       = "your_admin_user_name"
  cluster_admin_user_pwd        = "your_admin_password"
  workload_queue_name           = "your_queue_name"
  workload_queue_configurations = [
    {
      resource_name  = "cpu_limit"
      resource_value = "10"
    },
    {
      resource_name  = "memory"
      resource_value = "10"
    },
    {
      resource_name  = "tablespace"
      resource_value = "-1"
    },
    {
      resource_name  = "activestatements"
      resource_value = "-1"
    }
  ]

  user_name                          = "your_db_user"
  user_password                      = "your_db_user_password"
  workload_plan_name                 = "your_plan_name"
  workload_plan_stage_name           = "your_stage_name"
  workload_plan_stage_configurations = [
    {
      resource_name  = "cpu"
      resource_value = "1"
    },
    {
      resource_name  = "cpu_limit"
      resource_value = "0"
    },
    {
      resource_name  = "memory"
      resource_value = "10"
    },
    {
      resource_name  = "concurrency"
      resource_value = "10"
    },
    {
      resource_name  = "shortQueryConcurrencyNum"
      resource_value = "-1"
    }
  ]

  exception_rule_name           = "your_exception_rule"
  exception_rule_configurations = [
    {
      key   = "action"
      value = "abort"
    },
    {
      key   = "blocktime"
      value = "300"
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
* All resources will be created in the specified region
* Resource management is only supported for clusters with version `8.0.0` or later
* CPU exclusive limit is only supported for clusters with version `8.1.3` or later
* Resource management plan is only supported for clusters with version `8.1.0.100` or later

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.90.0 |
