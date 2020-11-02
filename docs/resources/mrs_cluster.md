---
subcategory: "MapReduce Service (MRS)"
---

# huaweicloud\_mrs\_cluster

Manages resource cluster within HuaweiCloud MRS.
This is an alternative to `huaweicloud_mrs_cluster_v1`

## Example Usage:  Creating a MRS cluster

```hcl
resource "huaweicloud_mrs_cluster" "cluster1" {
  cluster_name          = "mrs-cluster"
  cluster_version       = "MRS 1.8.10"
  cluster_type          = 0
  region                = "cn-north-1"
  available_zone_id     = "ae04cf9d61544df3806a3feeb401b204"
  billing_type          = 12
  master_node_num       = 2
  core_node_num         = 3
  master_node_size      = "c3.4xlarge.2.linux.bigdata"
  core_node_size        = "c3.xlarge.4.linux.bigdata"
  volume_type           = "SATA"
  volume_size           = 100
  safe_mode             = 0
  cluster_admin_secret  = var.admin_secret
  node_public_cert_name = var.keypair
  vpc_id                = "51edfb75-f9f0-4bbc-b4dc-21466b93f60d"
  subnet_id             = "1d7a8646-43ee-455a-a3ab-40da87a1304c"

  component_list {
    component_name = "Hadoop"
  }
  component_list {
    component_name = "Spark"
  }
  component_list {
    component_name = "Hive"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the msr cluster resource. If omitted, the provider-level region will work as default. Changing this creates a new msr cluster resource.

* `billing_type` - (Required) The value is 12, indicating on-demand payment.

* `region` - (Required) Cluster region information. Obtain the value from
    Regions and Endpoints.

* `cluster_name` - (Required) Cluster name, which is globally unique and contains
    only 1 to 64 letters, digits, hyphens (-), and underscores (_).

* `cluster_version` - (Optional) Version of the clusters. Currently, MRS 1.8.10, MRS 1.9.2
    and MRS 2.1.0 are supported.

* `cluster_type` - (Optional) Type of clusters 0: analysis cluster 1: streaming
    cluster The default value is 0.

* `master_node_num` - (Required) Number of Master nodes The value is 2.

* `master_node_size` - (Required) Best match based on several years of commissioning
    experience. MRS supports specifications of hosts, and host specifications are
    determined by CPUs, memory, and disks space. MRS supports instance specifications
	detailed in [MRS specifications](https://support.huaweicloud.com/en-us/api-mrs/mrs_01_9006.html)

* `core_node_num` - (Required) Number of Core nodes Value range: 3 to 100 A
    maximum of 100 Core nodes are supported by default. If more than 100 Core nodes
    are required, contact technical support engineers or invoke background APIs
    to modify the database.

* `core_node_size` - (Required) Instance specification of a Core node Configuration
    method of this parameter is identical to that of master_node_size.

* `available_zone_id` - (Required) ID of an available zone. Obtain the value
    from Regions and Endpoints.
	North China AZ1 (cn-north-1a): ae04cf9d61544df3806a3feeb401b204,
	North China AZ2 (cn-north-1b): d573142f24894ef3bd3664de068b44b0,
	East China AZ1 (cn-east-2a): 72d50cedc49846b9b42c21495f38d81c,
	East China AZ2 (cn-east-2b): 38b0f7a602344246bcb0da47b5d548e7,
	East China AZ3 (cn-east-2c): 5547fd6bf8f84bb5a7f9db062ad3d015,
	South China AZ1(cn-south-1a): 34f5ff4865cf4ed6b270f15382ebdec5,
	South China AZ2(cn-south-2b): 043c7e39ecb347a08dc8fcb6c35a274e,
	South China AZ3(cn-south-1c): af1687643e8c4ec1b34b688e4e3b8901,

* `vpc_id` - (Required) ID of the VPC where the subnet locates Obtain the VPC
    ID from the management console as follows: Register an account and log in to
    the management console. Click Virtual Private Cloud and select Virtual Private
    Cloud from the left list. On the Virtual Private Cloud page, obtain the VPC
    ID from the list.

* `subnet_id` - (Required) Subnet ID Obtain the subnet ID from the management
    console as follows: Register an account and log in to the management console.
    Click Virtual Private Cloud and select Virtual Private Cloud from the left list.
    On the Virtual Private Cloud page, obtain the subnet ID from the list.

* `volume_type` - (Required) Type of disks SATA and SSD are supported. SATA:
    common I/O SSD: super high-speed I/O

* `volume_size` - (Required) Data disk storage space of a Core node Users can
    add disks to expand storage capacity when creating a cluster. There are the
    following scenarios: Separation of data storage and computing: Data is stored
    in the OBS system. Costs of clusters are relatively low but computing performance
    is poor. The clusters can be deleted at any time. It is recommended when data
    computing is not frequently performed. Integration of data storage and computing:
    Data is stored in the HDFS system. Costs of clusters are relatively high but
    computing performance is good. The clusters cannot be deleted in a short term.
    It is recommended when data computing is frequently performed. Value range:
    100 GB to 32000 GB

* `node_public_cert_name` - (Required) Name of a key pair You can use a key
    to log in to the Master node in the cluster.

* `safe_mode` - (Required) MRS cluster running mode 0: common mode The value
    indicates that the Kerberos authentication is disabled. Users can use all functions
    provided by the cluster. 1: safe mode The value indicates that the Kerberos
    authentication is enabled. Common users cannot use the file management or job
    management functions of an MRS cluster and cannot view cluster resource usage
    or the job records of Hadoop and Spark. To use these functions, the users must
    obtain the relevant permissions from the MRS Manager administrator. The request
    has the cluster_admin_secret parameter only when safe_mode is set to 1.

* `cluster_admin_secret` - (Optional) Indicates the password of the MRS Manager
    administrator. The password for MRS 1.5.0: Must contain 6 to 32 characters.
    Must contain at least two types of the following: Lowercase letters Uppercase
    letters Digits Special characters of `~!@#$%^&*()-_=+\|[{}];:'",<.>/? Spaces
    Must be different from the username. Must be different from the username written
    in reverse order. The password for MRS 1.3.0: Must contain 8 to 64 characters.
    Must contain at least four types of the following: Lowercase letters Uppercase
    letters Digits Special characters of `~!@#$%^&*()-_=+\|[{}];:'",<.>/? Spaces
    Must be different from the username. Must be different from the username written
    in reverse order. This parameter needs to be configured only when safe_mode
    is set to 1.

* `log_collection` - (Optional) Indicates whether logs are collected when cluster
    installation fails. 0: not collected 1: collected The default value is 0. If
    log_collection is set to 1, OBS buckets will be created to collect the MRS logs.
    These buckets will be charged.

* `component_list` - (Required) Service component list.

* `add_jobs` - (Optional) You can submit a job when you create a cluster to
    save time and use MRS easily. Only one job can be added.

* `tags` - (Optional) The key/value pairs to associate with the cluster.

The `component_list` block supports:

* `component_name` - (Required) Component name Currently, Hadoop, Spark, HBase,
    Hive, Hue, Loader, Flume, Kafka and Storm are supported.


The `add_jobs` block supports:
* `job_type` - (Required) Job type 1: MapReduce 2: Spark 3: Hive Script 4: HiveQL
    (not supported currently) 5: DistCp, importing and exporting data (not supported
    in this API currently). 6: Spark Script 7: Spark SQL, submitting Spark SQL statements
    (not supported in this API currently). NOTE: Spark and Hive jobs can be added
    to only clusters including Spark and Hive components.

* `job_name` - (Required) Job name It contains only 1 to 64 letters, digits,
    hyphens (-), and underscores (_). NOTE: Identical job names are allowed but
    not recommended.

* `jar_path` - (Required) Path of the .jar file or .sql file for program execution
    The parameter must meet the following requirements: Contains a maximum of 1023
    characters, excluding special characters such as ;|&><'$. The address cannot
    be empty or full of spaces. Starts with / or s3a://. Spark Script must end with
    .sql; while MapReduce and Spark Jar must end with .jar. sql and jar are case-insensitive.

* `arguments` - (Optional) Key parameter for program execution The parameter
    is specified by the function of the user's program. MRS is only responsible
    for loading the parameter. The parameter contains a maximum of 2047 characters,
    excluding special characters such as ;|&>'<$, and can be empty.

* `input` - (Optional) Path for inputting data, which must start with / or s3a://.
    A correct OBS path is required. The parameter contains a maximum of 1023 characters,
    excluding special characters such as ;|&>'<$, and can be empty.

* `output` - (Optional) Path for outputting data, which must start with / or
    s3a://. A correct OBS path is required. If the path does not exist, the system
    automatically creates it. The parameter contains a maximum of 1023 characters,
    excluding special characters such as ;|&>'<$, and can be empty.

* `job_log` - (Optional) Path for storing job logs that record job running status.
    This path must start with / or s3a://. A correct OBS path is required. The parameter
    contains a maximum of 1023 characters, excluding special characters such as
    ;|&>'<$, and can be empty.

* `shutdown_cluster` - (Optional) Whether to delete the cluster after the jobs
    are complete true: Yes false: No

* `file_action` - (Optional) Data import and export import export

* `submit_job_once_cluster_run` - (Required) true: A job is submitted when a
    cluster is created. false: A job is submitted separately. The parameter is set
    to true in this example.

* `hql` - (Optional) HiveQL statement

* `hive_script_path` - (Optional) SQL program path This parameter is needed
    by Spark Script and Hive Script jobs only and must meet the following requirements:
    Contains a maximum of 1023 characters, excluding special characters such as
    ;|&><'$. The address cannot be empty or full of spaces. Starts with / or s3a://.
    Ends with .sql. sql is case-insensitive.

## Attributes Reference

The following attributes are exported:

* `billing_type` - See Argument Reference above.
* `data_center` - See Argument Reference above.
* `master_node_num` - See Argument Reference above.
* `master_node_size` - See Argument Reference above.
* `core_node_num` - See Argument Reference above.
* `core_node_size` - See Argument Reference above.
* `available_zone_id` - See Argument Reference above.
* `cluster_name` - See Argument Reference above.
* `vpc_id` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `cluster_version` - See Argument Reference above.
* `cluster_type` - See Argument Reference above.
* `volume_type` - See Argument Reference above.
* `volume_size` - See Argument Reference above.
* `node_public_cert_name` - See Argument Reference above.
* `safe_mode` - See Argument Reference above.
* `cluster_admin_secret` - See Argument Reference above.
* `log_collection` - See Argument Reference above.
* `component_list` - See Argument Reference below.
* `add_jobs` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `order_id` - Order ID for creating clusters.
* `cluster_id` - Cluster ID.
* `available_zone_name` - Name of an availability zone.
* `instance_id` - Instance ID.
* `hadoop_version` - Hadoop version.
* `master_node_ip` - IP address of a Master node.
* `externalIp` - Internal IP address.
* `private_ip_first` - Primary private IP address.
* `external_ip` - External IP address.
* `slave_security_groups_id` - Standby security group ID.
* `security_groups_id` - Security group ID.
* `external_alternate_ip` - Backup external IP address.
* `master_node_spec_id` - Specification ID of a Master node.
* `core_node_spec_id` - Specification ID of a Core node.
* `master_node_product_id` - Product ID of a Master node.
* `core_node_product_id` - Product ID of a Core node.
* `duration` - Cluster subscription duration.
* `vnc` - URI address for remote login of the elastic cloud server.
* `fee` - Cluster creation fee, which is automatically calculated.
* `deployment_id` - Deployment ID of a cluster.
* `cluster_state` - Cluster status Valid values include: existing history starting
    running terminated failed abnormal terminating rebooting shutdown frozen scaling-out
    scaling-in scaling-error.
* `tenant_id` - Project ID.
* `create_at` - Cluster creation time.
* `update_at` - Cluster update time.
* `error_info` - Error information.
* `charging_start_time` - Time when charging starts.
* `remark` - Remarks of a cluster.

The component_list attributes:
* `component_name` - (Required) Component name Currently, Hadoop, Spark, HBase,
    Hive, Hue, Loader, Flume, Kafka and Storm are supported.
