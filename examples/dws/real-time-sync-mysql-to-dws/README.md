# Sync MySQL data to DWS with DLI Flink

This example provides best practice code for using Terraform to sync MySQL data to a
GaussDB(DWS) cluster in real time with a DLI Flink OpenSource SQL job on HuaweiCloud.

## Architecture

1. Create a shared VPC/subnet and security group for RDS and DWS
2. Create an RDS MySQL instance and database
3. Create a DWS cluster
4. Create a DLI elastic resource pool (CIDR must differ from the VPC) and a general queue
5. Create an enhanced datasource connection and associate the resource pool

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC shared by RDS and DWS
* `vpc_cidr` - The CIDR block of the VPC. Must differ from the DLI elastic resource pool CIDR
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group shared by the RDS instance and DWS cluster
* `elastic_resource_pool_cidr` - The CIDR block of the DLI elastic resource pool. Must differ from the VPC CIDR
* `rds_instance_name` - The name of the RDS MySQL instance
* `rds_db_password` - The root password of the RDS MySQL instance
* `dws_cluster_name` - The name of the DWS cluster
* `dws_admin_user_pwd` - The administrator password of the DWS cluster
* `elastic_resource_pool_name` - The name of the DLI elastic resource pool
* `queue_name` - The name of the DLI general queue
* `datasource_connection_name` - The name of the DLI enhanced datasource connection

#### Optional Variables

* `availability_zone` - The availability zone. If empty, the first available zone is used (default: "")
* `enterprise_project_id` - The ID of the enterprise project (default: "")
* `subnet_cidr` - The CIDR block of the subnet. If empty, it is calculated from the VPC CIDR (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet. If empty, it is calculated from the subnet CIDR (default: "")
* `security_group_delete_default_rules` - Whether to delete the default rules of the security group (default: true)
* `dws_port` - The service port of the DWS cluster (default: 8000)
* `rds_db_port` - The database port of the RDS MySQL instance (default: 3306)
* `rds_flavor_id` - The flavor ID of the RDS instance. If empty, it is queried from huaweicloud_rds_flavors (default: "")
* `rds_db_version` - The MySQL version of the RDS instance (default: "5.7")
* `rds_instance_mode` - The instance mode used to query RDS flavors (default: "single")
* `rds_flavor_vcpus` - The vCPUs used to query RDS flavors (default: 2)
* `rds_volume_type` - The volume type of the RDS instance (default: "CLOUDSSD")
* `rds_volume_size` - The volume size of the RDS instance in GB (default: 40)
* `dws_node_type` - The flavor of the DWS cluster node. If empty, it is queried from huaweicloud_dws_flavors (default: "")
* `dws_version` - The version of the DWS cluster. If empty, it is queried from huaweicloud_dws_flavors (default: "")
* `dws_flavor_vcpus` - The vCPUs used to query DWS flavors (default: 4)
* `dws_flavor_memory` - The memory used to query DWS flavors (default: 32)
* `dws_datastore_type` - The datastore type of the DWS cluster (default: "dws")
* `dws_number_of_node` - The number of nodes in the DWS cluster (default: 3)
* `dws_number_of_cn` - The number of CN nodes in the DWS cluster (default: 3)
* `dws_admin_user_name` - The administrator username of the DWS cluster (default: "dbadmin")
* `dws_volume_type` - The volume type of the DWS cluster (default: "SSD")
* `dws_volume_capacity` - The volume capacity of the DWS cluster in GB (default: "100")
* `elastic_resource_pool_description` - The description of the DLI elastic resource pool (default: "")
* `elastic_resource_pool_min_cu` - The minimum number of CUs for the DLI elastic resource pool (default: 16)
* `elastic_resource_pool_max_cu` - The maximum number of CUs for the DLI elastic resource pool (default: 64)
* `elastic_resource_pool_label` - The label of the DLI elastic resource pool (default: `{ spec = "basic" }`)
* `queue_cu_count` - The CU count of the DLI queue (default: 16)
* `queue_description` - The description of the DLI queue (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                   = "your_vpc_name"
  vpc_cidr                   = "192.168.0.0/16"
  subnet_name                = "your_subnet_name"
  security_group_name        = "your_security_group_name"
  elastic_resource_pool_cidr = "172.16.0.0/18"
  rds_instance_name          = "your_rds_instance_name"
  rds_db_password            = "your_rds_password"
  dws_cluster_name           = "your_dws_cluster_name"
  dws_admin_user_pwd         = "your_dws_password"
  elastic_resource_pool_name = "your_dli_pool_name"
  queue_name                 = "your_dli_queue_name"
  datasource_connection_name = "your_datasource_connection_name"
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

## Configuration Details

### Create Tables on RDS and DWS

After the infrastructure is created, prepare the source and target tables.

#### RDS MySQL source table

Connect to the RDS instance (DAS console or a host that can reach the instance), create the database if needed,
then create the source table and insert sample data:

```sql
CREATE DATABASE IF NOT EXISTS mys_data;

CREATE TABLE mys_data.mys_order (
  order_id      VARCHAR(12),
  order_channel VARCHAR(32),
  order_time    DATETIME,
  cust_code     VARCHAR(6),
  pay_amount    DOUBLE,
  real_pay      DOUBLE,
  PRIMARY KEY (order_id)
);

INSERT INTO mys_data.mys_order VALUES
  ('202306270001', 'webShop', '2023-06-27 10:00:00', 'CUST1', 1000, 1000),
  ('202306270002', 'webShop', '2023-06-27 11:00:00', 'CUST2', 5000, 5000);

SELECT * FROM mys_data.mys_order;
```

#### DWS target schema and table

Log in to the DWS cluster (console **Login** or a SQL client), then create the schema and target table:

```sql
CREATE SCHEMA dws_data;

CREATE TABLE dws_data.dws_order (
  order_id      VARCHAR(12),
  order_channel VARCHAR(32),
  order_time    TIMESTAMP,
  cust_code     VARCHAR(6),
  pay_amount    DOUBLE PRECISION,
  real_pay      DOUBLE PRECISION
);

SELECT * FROM dws_data.dws_order;
```

The query should return an empty result before the Flink job starts syncing data.

### Prepare the DWS Connector JAR

Download a Flink **1.15**-compatible `dws-connector-flink` JAR from
[Maven Repository - com.huaweicloud.dws](https://mvnrepository.com/artifact/com.huaweicloud.dws)
(for example `dws-connector-flink-sql-1.15-*.jar`), then upload it to an OBS bucket in the same region as DLI.
The Flink job can reference the JAR OBS path as the UDF Jar.

### Create the DLI Flink OpenSource SQL Job

Create a Flink OpenSource SQL job that reads MySQL with `mysql-cdc` and writes to DWS with the `gaussdb` connector.

Recommended job settings:

* **Queue**: the general queue created in this example
* **Flink version**: `1.15` or later
* **UDF Jar**: the `dws-connector-flink` JAR uploaded to OBS above
* **Agency**: the DLI agency authorized to access OBS
* **OBS bucket**: enable job logs and checkpoint as needed
* **Checkpoint**: enable for better reliability

Example Flink SQL (replace the RDS private IP, DWS private IP, and passwords):

```sql
CREATE TABLE mys_order (
  order_id STRING,
  order_channel STRING,
  order_time TIMESTAMP,
  cust_code STRING,
  pay_amount DOUBLE,
  real_pay DOUBLE,
  PRIMARY KEY (order_id) NOT ENFORCED
) WITH (
  'connector' = 'mysql-cdc',
  'hostname' = '<RDS_PRIVATE_IP>',
  'port' = '3306',
  'username' = 'root',
  'password' = '<RDS_ROOT_PASSWORD>',
  'database-name' = 'mys_data',
  'table-name' = 'mys_order'
);

CREATE TABLE dws_order (
  order_id STRING,
  order_channel STRING,
  order_time TIMESTAMP,
  cust_code STRING,
  pay_amount DOUBLE,
  real_pay DOUBLE,
  PRIMARY KEY (order_id) NOT ENFORCED
) WITH (
  'connector' = 'gaussdb',
  'driver' = 'com.huawei.gauss200.jdbc.Driver',
  'url' = 'jdbc:gaussdb://<DWS_PRIVATE_IP>:8000/gaussdb',
  'table-name' = 'dws_data.dws_order',
  'username' = 'dbadmin',
  'password' = '<DWS_DBADMIN_PASSWORD>',
  'write.mode' = 'insert'
);

INSERT INTO dws_order SELECT * FROM mys_order;
```

Common parameters:

* `connector`: source uses `mysql-cdc`; sink uses `gaussdb`
* `driver`: fixed as `com.huawei.gauss200.jdbc.Driver` for DWS
* `write.mode`: supports `copy`, `insert`, and `upsert`

After formatting and saving the SQL, start the job. When the status becomes **running**, verify data in DWS:

```sql
SELECT * FROM dws_data.dws_order;
```

You can insert more rows into MySQL and query DWS again to confirm near-real-time sync.

For production, prefer DLI datasource authentication (`pwd_auth_name`) instead of putting database passwords
directly in the Flink SQL script.

## Note

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region
* DWS and RDS must be in the same region and VPC
* The DLI elastic resource pool CIDR must not overlap with the VPC CIDR of RDS/DWS
* After apply, test datasource connectivity from the DLI queue to the RDS and DWS private IPs in the console

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.75.5 |
