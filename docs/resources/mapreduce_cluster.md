---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_cluster"
description: ""
---

# huaweicloud_mapreduce_cluster

Manages a cluster resource within HuaweiCloud MRS.

## Example Usage

### Create an analysis cluster

```hcl
data "huaweicloud_availability_zones" "test" {}

variable "cluster_name" {}
variable "password" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = var.cluster_name
  version            = "MRS 1.9.2"
  type               = "ANALYSIS"
  component_list     = ["Hadoop", "Hive", "Tez"]
  manager_admin_pass = var.password
  node_admin_pass    = var.password
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_task_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a stream cluster

```hcl
data "huaweicloud_availability_zones" "test" {}

variable "cluster_name" {}
variable "password" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = var.cluster_name
  type               = "STREAMING"
  version            = "MRS 1.9.2"
  manager_admin_pass = var.password
  node_admin_pass    = var.password
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  component_list     = ["Storm"]

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a hybrid cluster

```hcl
data "huaweicloud_availability_zones" "test" {}

variable "cluster_name" {}
variable "password" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = var.cluster_name
  version            = "MRS 3.1.5"
  type               = "MIXED"
  component_list     = ["Hadoop", "ZooKeeper", "Ranger", "Tez", "Spark2x", "Hive", "Kafka", "Flume"]
  manager_admin_pass = var.password
  node_admin_pass    = var.password
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id

  master_nodes {
    flavor            = "ac7.4xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "ac7.4xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "ac7.4xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
  }
  analysis_task_nodes {
    flavor            = "ac7.4xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a custom cluster

```hcl
data "huaweicloud_availability_zones" "test" {}

variable "cluster_name" {}
variable "password" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = var.cluster_name
  version            = "MRS 3.1.0"
  type               = "CUSTOM"
  safe_mode          = true
  manager_admin_pass = var.password
  node_admin_pass    = var.password
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  template_id        = "mgmt_control_combined_v4"
  component_list     = ["DBService", "Hadoop", "ZooKeeper", "Ranger"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 4
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }
}

```

### Create an analysis cluster and bind public IP

```hcl
data "huaweicloud_availability_zones" "test" {}

variable "cluster_name" {}
variable "password" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "public_ip" {}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = var.cluster_name
  version            = "MRS 1.9.2"
  type               = "ANALYSIS"
  component_list     = ["Hadoop", "Hive", "Tez"]
  manager_admin_pass = var.password
  node_admin_pass    = var.password
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  public_ip          = var.public_ip

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_task_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a prePaid stream cluster

```hcl
data "huaweicloud_availability_zones" "test" {}

variable "cluster_name" {}
variable "password" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "period" {}
variable "period_unit" {}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = var.cluster_name
  type               = "STREAMING"
  version            = "MRS 1.9.2"
  manager_admin_pass = var.password
  node_admin_pass    = var.password
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  component_list     = ["Storm"]
  charging_mode      = "prePaid"
  period             = var.period
  period_unit        = var.period_unit

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  streaming_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the MapReduce cluster resource. If omitted, the
  provider-level region will be used. Changing this will create a new MapReduce cluster resource.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone in which to create the cluster.
  Please following [reference](https://developer.huaweicloud.com/intl/en-us/endpoint?all)
  Changing this will create a new MapReduce cluster resource.

* `name` - (Required, String) Specifies the name of the MapReduce cluster. The name can contain 2 to 64
  characters, which may consist of letters, digits, underscores (_) and hyphens (-).

* `version` - (Required, String, ForceNew) Specifies the MapReduce cluster version. The valid values are `MRS 1.9.2`
  , `MRS 3.0.5` and `MRS 3.1.0`. Changing this will create a new MapReduce cluster resource.

* `component_list` - (Required, List, ForceNew) Specifies the list of component names. For the components supported by
  the cluster, please following [reference](https://support.huaweicloud.com/intl/en-us/productdesc-mrs/mrs_08_0005.html)
  Changing this will create a new MapReduce cluster resource.

* `master_nodes` - (Required, List, ForceNew) Specifies the informations about master nodes in the MapReduce cluster.
  The [master_nodes](#v2_mapreduce_cluster_master_nodes) structure is documented below.
  Changing this will create a new MapReduce cluster resource.

* `manager_admin_pass` - (Required, String, ForceNew) Specifies the administrator password, which is used to log in to
  the cluster management page. The password can contain 8 to 26 characters and cannot be the username or the username
  spelled backwards. The password must contain lowercase letters, uppercase letters, digits, spaces and the special
  characters: `!?,.:-_{}[]@$^+=/`. Changing this will create a new MapReduce cluster resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC which bound to the MapReduce cluster. Changing
  this will create a new MapReduce cluster resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of a subnet which bound to the MapReduce cluster.
  Changing this will create a new MapReduce cluster resource.

* `type` - (Optional, String, ForceNew) Specifies the type of the MapReduce cluster. The valid values are **ANALYSIS**,
  **STREAMING** and **MIXED**, defaults to **ANALYSIS**. Changing this will create a new MapReduce cluster resource.

* `enterprise_project_id` - (Optional, String) Specifies a unique ID in UUID format of enterprise project.

* `public_ip` - (Optional, String, ForceNew) Specifies the EIP address which bound to the MapReduce cluster.
The EIP must have been created and must be in the same region as the cluster.
 Changing this will create a new MapReduce cluster resource.

* `eip_id` - (Optional, String, ForceNew) Specifies the EIP ID which bound to the MapReduce cluster.
The EIP must have been created and must be in the same region as the cluster.
 Changing this will create a new MapReduce cluster resource.

* `log_collection` - (Optional, Bool, ForceNew) Specifies whether logs are collected when cluster installation fails.
  Defaults to true. If `log_collection` set true, the OBS buckets will be created and only used to collect logs that
  record MapReduce cluster creation failures. Changing this will create a new MapReduce cluster resource.

* `node_admin_pass` - (Optional, String, ForceNew) Specifies the administrator password, which is used to log in to the
  each nodes(/ECSs). The password can contain 8 to 26 characters and cannot be the username or the username spelled
  backwards. The password must contain lowercase letters, uppercase letters, digits, spaces and the special
  characters: `!?,.:-_{}[]@$^+=/`. Changing this will create a new MapReduce cluster resource. This parameter
  and `node_key_pair` are alternative.

* `node_key_pair` - (Optional, String, ForceNew) Specifies the name of a key pair, which is used to log in to the each
  nodes(/ECSs). Changing this will create a new MapReduce cluster resource.

* `safe_mode` - (Optional, Bool, ForceNew) Specifies whether the running mode of the MapReduce cluster is secure,
  defaults to **true**. The options are as follows:
  + **true**: enable Kerberos authentication.
  + **false**: disable Kerberos authentication. Changing this will create a new MapReduce cluster resource.

* `security_group_ids` - (Optional, List, ForceNew) Specifies an array of one or more security group ID to attach to the
  MapReduce cluster. If using the specified security group, the group need to open the specified port (9022) rules.

* `template_id` - (Optional, String, ForceNew) Specifies the template used for node deployment when the cluster type is
  **CUSTOM**. The options are as follows:
  + **mgmt_control_combined_v2**: template for jointly deploying the management and control nodes. The management and
  control roles are co-deployed on the Master node, and data instances are deployed in the same node group. This
  deployment mode applies to scenarios where the number of control nodes is less than 100, reducing costs.
  + **mgmt_control_separated_v2**: The management and control roles are deployed on different master nodes, and data
  instances are deployed in the same node group. This deployment mode is applicable to a cluster with 100 to 500 nodes
  and delivers better performance in high-concurrency load scenarios.
  + **mgmt_control_data_separated_v2**: The management role and control role are deployed on different Master nodes,
  and data instances are deployed in different node groups. This deployment mode is applicable to a cluster with more
  than 500 nodes. Components can be deployed separately, which can be used for a larger cluster scale.

* `analysis_core_nodes` - (Optional, List, ForceNew) Specifies the informations about analysis core nodes in the
  MapReduce cluster.  
  The [analysis_core_nodes](#v2_mapreduce_cluster_core_nodes) structure is documented below.  
  Changing this will create a new MapReduce cluster resource.

* `streaming_core_nodes` - (Optional, List, ForceNew) Specifies the informations about streaming core nodes in the
  MapReduce cluster.  
  The [streaming_core_nodes](#v2_mapreduce_cluster_core_nodes) structure is documented below.  
  Changing this will create a new MapReduce cluster resource.

* `analysis_task_nodes` - (Optional, List, ForceNew) Specifies the informations about analysis task nodes in the
  MapReduce cluster.  
  The [analysis_task_nodes](#v2_mapreduce_cluster_task_nodes) structure is documented below.  
  Changing this will create a new MapReduce cluster resource.

* `streaming_task_nodes` - (Optional, List, ForceNew) Specifies the informations about streaming task nodes in the
  MapReduce cluster.  
  The [streaming_task_nodes](#v2_mapreduce_cluster_task_nodes) structure is documented below.  
  Changing this will create a new MapReduce cluster resource.

* `custom_nodes` - (Optional, List, ForceNew) Specifies the informations about custom nodes in the MapReduce cluster.  
  The [custom_nodes](#v2_mapreduce_cluster_custom_nodes) structure is documented below.  
  Changing this will create a new MapReduce cluster resource.

* `component_configs` - (Optional, List, ForceNew) Specifies the component configurations of the cluster.
  The [object](#component_configurations) structure is documented below.
  Changing this will create a new MapReduce cluster resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the cluster.

* `external_datasources` - (Optional, List, ForceNew) Specifies the external datasource configurations of the cluster.
  The [object](#ExternalDatasources) structure is documented below.
  Changing this will create a new MapReduce cluster resource.

* `bootstrap_scripts` - (Optional, List, ForceNew) Specifies the bootstrap action scripts.
  Bootstrap action scripts will be executed on specified cluster nodes before or after big data components are
  started. You can execute bootstrap actions to install third-party software, modify the cluster running environment,
  and customize cluster configuration. [Learn more](https://support.huaweicloud.com/intl/en-us/usermanual-mrs/mrs_01_0414.html)

  The [object](#BootstrapScripts) structure is documented below.
  Changing this will create a new MapReduce cluster resource.

* `smn_notify` - (Optional, List, ForceNew) Specifies the alarm configuration of the cluster. The [smn_notify](#SMNNotify)
  structure is documented below.
  Changing this will create a new MapReduce cluster resource.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the cluster.  
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the cluster.  
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the cluster.  
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.  
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.  
  This parameter is mandatory if `charging_mode` is set to **prePaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled, defaults to **false**.  
  This parameter is available if `charging_mode` is set to **prePaid**.  
  The valid values are **true** and **false**.  

<a name="v2_mapreduce_cluster_master_nodes"></a>
The `master_nodes` block supports:

* `flavor` - (Required, String, ForceNew) Specifies the instance specifications for each nodes in node group.
  Changing this will create a new MapReduce cluster resource.

* `node_number` - (Required, Int, ForceNew) Specifies the number of nodes for the node group.  
  Changing this will create a new MapReduce cluster resource.

* `root_volume_type` - (Required, String, ForceNew) Specifies the system disk flavor of the nodes. Changing this will
  create a new MapReduce cluster resource.

* `root_volume_size` - (Required, Int, ForceNew) Specifies the system disk size of the nodes. Changing this will create
  a new MapReduce cluster resource.

* `data_volume_count` - (Required, Int, ForceNew) Specifies the data disk number of the nodes.  
  The valid value is `1`.  
  Changing this will create a new MapReduce cluster resource.

* `data_volume_type` - (Optional, String, ForceNew) Specifies the data disk flavor of the nodes.
  Required if `data_volume_count` is greater than zero. Changing this will create a new MapReduce cluster resource.
   The following disk types are supported:
  + **SATA**: common I/O disk.
  + **SAS**: high I/O disk.
  + **SSD**: ultra-high I/O disk.

* `data_volume_size` - (Optional, Int, ForceNew) Specifies the data disk size of the nodes,in GB. The value range is 10
  to 32768. Required if `data_volume_count` is greater than zero. Changing this will create a new MapReduce
  cluster resource.

* `assigned_roles` - (Optional, List, ForceNew) Specifies the roles deployed in a node group.This argument is mandatory
 when the cluster type is **CUSTOM**. Each character string represents a role expression.

  **Role expression definition:**

   + If the role is deployed on all nodes in the node group, set this parameter to role_name, for example: `DataNode`.
   + If the role is deployed on a specified subscript node in the node group: role_name:index1,index2..., indexN,
 for example: `DataNode:1,2`. The subscript starts from 1.
   + Some roles support multi-instance deployment (that is, multiple instances of the same role are deployed on a node):
  role_name[instance_count], for example: `EsNode[9]`.
  
  [For details about components](https://support.huaweicloud.com/intl/en-us/productdesc-mrs/mrs_08_0005.html)

  [Mapping between roles and components](https://support.huaweicloud.com/intl/en-us/api-mrs/mrs_02_0106.html)

  -> `DBService` is a basic component of a cluster. Components such as Hive, Hue, Oozie, Loader, and Redis, and Loader
   store their metadata in DBService, and provide the metadata backup and restoration functions by using DBService.
  
* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the cluster.  
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the cluster.  
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the cluster.  
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.  
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.  
  This parameter is mandatory if `charging_mode` is set to **prePaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled, defaults to **false**.  
  This parameter is available if `charging_mode` is set to **prePaid**.  
  The valid values are **true** and **false**.  

  -> The `period_unit` must be used together with the `period` parameter.
  
<a name="v2_mapreduce_cluster_core_nodes"></a>
The `analysis_core_nodes` and `streaming_core_nodes` blocks support:

* `flavor` - (Required, String, ForceNew) Specifies the instance specifications for each nodes in node group.
  Changing this will create a new MapReduce cluster resource.

* `node_number` - (Required, Int) Specifies the number of nodes for the node group.

  -> The number of nodes after scaling cannot be less than the number of nodes originally created.

* `root_volume_type` - (Required, String, ForceNew) Specifies the system disk flavor of the nodes. Changing this will
  create a new MapReduce cluster resource.

* `root_volume_size` - (Required, Int, ForceNew) Specifies the system disk size of the nodes. Changing this will create
  a new MapReduce cluster resource.

* `data_volume_count` - (Required, Int, ForceNew) Specifies the data disk number of the nodes. The number configuration
  of each node are as follows:
  + **analysis_core_nodes**: minimum is one and the maximum is subject to the configuration of the corresponding flavor.
  + **streaming_core_nodes**: minimum is one and the maximum is subject to the configuration of the corresponding flavor.
  + **analysis_task_nodes**: minimum is one and the maximum is subject to the configuration of the corresponding flavor.
  + **streaming_task_nodes**: minimum is one and the maximum is subject to the configuration of the corresponding flavor.

  Changing this will create a new MapReduce cluster resource.
  
* `data_volume_type` - (Optional, String, ForceNew) Specifies the data disk flavor of the nodes.
  Required if `data_volume_count` is greater than zero. Changing this will create a new MapReduce cluster resource.
   The following disk types are supported:
  + **SATA**: common I/O disk.
  + **SAS**: high I/O disk.
  + **SSD**: ultra-high I/O disk.

* `data_volume_size` - (Optional, Int, ForceNew) Specifies the data disk size of the nodes,in GB. The value range is 10
  to 32768. Required if `data_volume_count` is greater than zero. Changing this will create a new MapReduce
  cluster resource.

* `assigned_roles` - (Optional, List, ForceNew) Specifies the roles deployed in a node group.This argument is mandatory
 when the cluster type is **CUSTOM**. Each character string represents a role expression.

  **Role expression definition:**

   + If the role is deployed on all nodes in the node group, set this parameter to role_name, for example: `DataNode`.
   + If the role is deployed on a specified subscript node in the node group: role_name:index1,index2..., indexN,
 for example: `DataNode:1,2`. The subscript starts from 1.
   + Some roles support multi-instance deployment (that is, multiple instances of the same role are deployed on a node):
  role_name[instance_count], for example: `EsNode[9]`.
  
  [For details about components](https://support.huaweicloud.com/intl/en-us/productdesc-mrs/mrs_08_0005.html)

  [Mapping between roles and components](https://support.huaweicloud.com/intl/en-us/api-mrs/mrs_02_0106.html)

  -> `DBService` is a basic component of a cluster. Components such as Hive, Hue, Oozie, Loader, and Redis, and Loader
   store their metadata in DBService, and provide the metadata backup and restoration functions by using DBService.

  + `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the cluster.  
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the cluster.  
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the cluster.  
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.  
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.  
  This parameter is mandatory if `charging_mode` is set to **prePaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled, defaults to **false**.  
  This parameter is available if `charging_mode` is set to **prePaid**.  
  The valid values are **true** and **false**.  

  -> The `period_unit` must be used together with the `period` parameter.

<a name="v2_mapreduce_cluster_task_nodes"></a>
The `analysis_task_nodes` and `streaming_task_nodes` blocks support:

* `flavor` - (Required, String, ForceNew) Specifies the instance specifications for each nodes in node group.
  Changing this will create a new MapReduce cluster resource.

* `node_number` - (Required, Int) Specifies the number of nodes for the node group.

  -> The number of nodes after scaling cannot be less than the number of nodes originally created.

* `root_volume_type` - (Required, String, ForceNew) Specifies the system disk flavor of the nodes. Changing this will
  create a new MapReduce cluster resource.

* `root_volume_size` - (Required, Int, ForceNew) Specifies the system disk size of the nodes. Changing this will create
  a new MapReduce cluster resource.

* `data_volume_count` - (Required, Int, ForceNew) Specifies the data disk number of the nodes. The number configuration
  of each node are as follows:
  + **analysis_core_nodes**: minimum is one and the maximum is subject to the configuration of the corresponding flavor.
  + **streaming_core_nodes**: minimum is one and the maximum is subject to the configuration of the corresponding flavor.
  + **analysis_task_nodes**: minimum is one and the maximum is subject to the configuration of the corresponding flavor.
  + **streaming_task_nodes**: minimum is one and the maximum is subject to the configuration of the corresponding flavor.

  Changing this will create a new MapReduce cluster resource.
  
* `data_volume_type` - (Optional, String, ForceNew) Specifies the data disk flavor of the nodes.
  Required if `data_volume_count` is greater than zero. Changing this will create a new MapReduce cluster resource.
   The following disk types are supported:
  + **SATA**: common I/O disk.
  + **SAS**: high I/O disk.
  + **SSD**: ultra-high I/O disk.

* `data_volume_size` - (Optional, Int, ForceNew) Specifies the data disk size of the nodes,in GB. The value range is 10
  to 32768. Required if `data_volume_count` is greater than zero. Changing this will create a new MapReduce
  cluster resource.

* `assigned_roles` - (Optional, List, ForceNew) Specifies the roles deployed in a node group.This argument is mandatory
 when the cluster type is **CUSTOM**. Each character string represents a role expression.

  **Role expression definition:**

   + If the role is deployed on all nodes in the node group, set this parameter to role_name, for example: `DataNode`.
   + If the role is deployed on a specified subscript node in the node group: role_name:index1,index2..., indexN,
 for example: `DataNode:1,2`. The subscript starts from 1.
   + Some roles support multi-instance deployment (that is, multiple instances of the same role are deployed on a node):
  role_name[instance_count], for example: `EsNode[9]`.
  
  [For details about components](https://support.huaweicloud.com/intl/en-us/productdesc-mrs/mrs_08_0005.html)

  [Mapping between roles and components](https://support.huaweicloud.com/intl/en-us/api-mrs/mrs_02_0106.html)

  -> `DBService` is a basic component of a cluster. Components such as Hive, Hue, Oozie, Loader, and Redis, and Loader
   store their metadata in DBService, and provide the metadata backup and restoration functions by using DBService.

<a name="v2_mapreduce_cluster_custom_nodes"></a>
The `custom_nodes` block supports:

* `group_name` - (Required, String, ForceNew) Specifies the name of nodes for the node group.  
  Changing this will create a new MapReduce cluster resource.

* `flavor` - (Required, String, ForceNew) Specifies the instance specifications for each nodes in node group.
  Changing this will create a new MapReduce cluster resource.

* `node_number` - (Required, Int, ForceNew) Specifies the number of nodes for the node group.  
  Changing this will create a new MapReduce cluster resource.

* `root_volume_type` - (Required, String, ForceNew) Specifies the system disk flavor of the nodes. Changing this will
  create a new MapReduce cluster resource.

* `root_volume_size` - (Required, Int, ForceNew) Specifies the system disk size of the nodes. Changing this will create
  a new MapReduce cluster resource.

* `data_volume_count` - (Required, Int, ForceNew) Specifies the data disk number of the nodes.  
  The valid value is `1`.  
  Changing this will create a new MapReduce cluster resource.

* `data_volume_type` - (Optional, String, ForceNew) Specifies the data disk flavor of the nodes.
  Required if `data_volume_count` is greater than zero. Changing this will create a new MapReduce cluster resource.
   The following disk types are supported:
  + **SATA**: common I/O disk.
  + **SAS**: high I/O disk.
  + **SSD**: ultra-high I/O disk.

* `data_volume_size` - (Optional, Int, ForceNew) Specifies the data disk size of the nodes,in GB. The value range is 10
  to 32768. Required if `data_volume_count` is greater than zero. Changing this will create a new MapReduce
  cluster resource.

* `assigned_roles` - (Optional, List, ForceNew) Specifies the roles deployed in a node group.This argument is mandatory
 when the cluster type is **CUSTOM**. Each character string represents a role expression.

  **Role expression definition:**

   + If the role is deployed on all nodes in the node group, set this parameter to role_name, for example: `DataNode`.
   + If the role is deployed on a specified subscript node in the node group: role_name:index1,index2..., indexN,
 for example: `DataNode:1,2`. The subscript starts from 1.
   + Some roles support multi-instance deployment (that is, multiple instances of the same role are deployed on a node):
  role_name[instance_count], for example: `EsNode[9]`.
  
  [For details about components](https://support.huaweicloud.com/intl/en-us/productdesc-mrs/mrs_08_0005.html)

  [Mapping between roles and components](https://support.huaweicloud.com/intl/en-us/api-mrs/mrs_02_0106.html)

  -> `DBService` is a basic component of a cluster. Components such as Hive, Hue, Oozie, Loader, and Redis, and Loader
   store their metadata in DBService, and provide the metadata backup and restoration functions by using DBService.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the cluster.  
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the cluster.  
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the cluster.  
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.  
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.  
  This parameter is mandatory if `charging_mode` is set to **prePaid**.  
  Changing this parameter will create a new MapReduce cluster resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled, defaults to **false**.  
  This parameter is available if `charging_mode` is set to **prePaid**.  
  The valid values are **true** and **false**.  

  -> The `period_unit` must be used together with the `period` parameter.

<a name="component_configurations"></a>
The `component_configs` block supports:

* `name` - (Required, String, ForceNew) Specifies the component name of the cluster which has installed.
  Changing this will create a new MapReduce cluster resource.

* `configs` - (Required, List, ForceNew) Specifies the configuration of component installed.
  The [object](#component_configuration) structure is documented below.

<a name="component_configuration"></a>
The `configs` block supports:

* `key` - (Required, String, ForceNew) Specifies the configuration item key of component installed.
  Changing this will create a new MapReduce cluster resource.

* `value` - (Required, String, ForceNew) Specifies the configuration item value of component installed.
  Changing this will create a new MapReduce cluster resource.

* `config_file_name` - (Required, String, ForceNew) Specifies the configuration file name of component installed.
  Changing this will create a new MapReduce cluster resource.

<a name="ExternalDatasources"></a>
The `external_datasources` block supports:

* `component_name` - (Required, String, ForceNew) Specifies the component name. The valid values are `Hive` and `Ranger`.
  Changing this will create a new MapReduce cluster resource.

* `role_type` - (Required, String, ForceNew) Specifies the component role type.
  The options are as follows:
    + **hive_metastore**: Hive Metastore role.
    + **ranger_data**: Ranger role.
  
  Changing this will create a new MapReduce cluster resource.

* `source_type` - (Required, String, ForceNew) Specifies the data connection type.
  The options are as follows:
    + **LOCAL_DB**: Local metadata.
    + **RDS_POSTGRES**: RDS PostgreSQL database.
    + **RDS_MYSQL**: RDS MySQL database.
    + **gaussdb-mysql**: GaussDB(for MySQL).
  
  Changing this will create a new MapReduce cluster resource.

* `data_connection_id` - (Optional, String, ForceNew) Specifies the data connection ID.
  This parameter is mandatory if `source_type` is not **LOCAL_DB**.
  Changing this will create a new MapReduce cluster resource.

<a name="BootstrapScripts"></a>
The `bootstrap_scripts` block supports:

* `name` - (Required, String, ForceNew) Specifies the name of a bootstrap action script.
  Changing this will create a new MapReduce cluster resource.

* `uri` - (Required, String, ForceNew) Specifies the path of a bootstrap action script.
  Set this parameter to an OBS bucket path or a local VM path.
    + **OBS bucket path**: The path of an OBS file system starts with *s3a://* or *obs://* and end with *.sh*.
    + **Local VM path**: The script path must start with a slash (/) and end with *.sh*.
  
  Changing this will create a new MapReduce cluster resource.

* `nodes` - (Required, List, ForceNew) Specifies names of the node group where the bootstrap action script is executed.

* `fail_action` - (Required, String, ForceNew) Specifies the action after the bootstrap action script fails to be
  executed.
  The options are as follows:
    + **continue**: Continue to execute subsequent scripts.
    + **errorout**: Stop the action.
  
  The default value is **errorout**, indicating that the action is stopped.
  Changing this will create a new MapReduce cluster resource.

  -> You are advised to set this parameter to continue in the commissioning phase so that the cluster can
     continue to be installed and started no matter whether the bootstrap action is successful.

* `parameters` - (Optional, String, ForceNew) Specifies bootstrap action script parameters.

* `active_master` - (Optional, Bool, ForceNew) Specifies whether the bootstrap action script runs only on active master
  nodes. The default value is **false**, indicating that the bootstrap action script can run on all master nodes.

* `before_component_start` - (Optional, Bool, ForceNew) Specifies whether the bootstrap action script is executed
  before component start.
  The options are as follows:
    + **false**: After component start. The default value is **false**.
    + **true**: Before component start.
  
  Changing this will create a new MapReduce cluster resource.

* `execute_need_sudo_root` - (Optional, Bool, ForceNew) Specifies whether the bootstrap action script involves root user
  operations.
  Changing this will create a new MapReduce cluster resource.

<a name="SMNNotify"></a>
The `smn_notify` block supports:

* `topic_urn` - (Required, String, ForceNew) Specifies the Uniform Resource Name (URN) of the topic.
  Changing this will create a new MapReduce cluster resource.

* `subscription_name` - (Required, String, ForceNew) Specifies the subscription rule name.
  Changing this will create a new MapReduce cluster resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The cluster ID in UUID format.
* `total_node_number` - The total number of nodes deployed in the cluster.
* `master_node_ip` - The IP address of the master node.
* `private_ip` - The preferred private IP address of the master node.
* `status` - The cluster state, which include: running, frozen, abnormal and failed.
* `create_time` - The cluster creation time, in RFC-3339 format.
* `update_time` - The cluster update time, in RFC-3339 format.
* `charging_start_time` - The charging start time which is the start time of billing, in RFC-3339 format.
* `node` - all the nodes attributes: master_nodes/analysis_core_nodes/streaming_core_nodes/analysis_task_nodes
  /streaming_task_nodes.
  + `host_ips` - The host list of this nodes group in the cluster.
* `bootstrap_scripts/start_time` - The execution time of one bootstrap action script, in RFC-3339 format.
* `bootstrap_scripts/state` - The status of one bootstrap action script.
    The valid value are **PENDING**, **IN_PROGRESS**, **SUCCESS**, and **FAILURE**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 180 minutes.
* `delete` - Default is 40 minutes.

## Import

Clusters can be imported by their `id`. For example,

```bash
terraform import huaweicloud_mapreduce_cluster.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`manager_admin_pass`, `node_admin_pass`,`template_id`, `assigned_roles`, `external_datasources`, `component_configs`,
`smn_notify`, `charging_mode`, `period`, `period_unit`, `auto_renew`, `streaming_core_nodes.0.charging_mode`,
`streaming_core_nodes.0.period`, `streaming_core_nodes.0.period_unit`, `streaming_core_nodes.0.auto_renew`,
`analysis_core_nodes.0.charging_mode`, `analysis_core_nodes.0.period`, `analysis_core_nodes.0.period_unit`,
`analysis_core_nodes.0.auto_renew`, `custom_nodes.0.charging_mode`, `custom_nodes.0.period`, `custom_nodes.0.period_unit`
and `custom_nodes.0.auto_renew`. It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition
should be updated to align with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_mapreduce_cluster" "test" {
    ...

  lifecycle {
    ignore_changes = [
      manager_admin_pass, node_admin_pass, template_id, assigned_roles, external_datasources, component_configs, smn_notify,
      charging_mode, period, period_unit, auto_renew, master_nodes.0.charging_mode, master_nodes.0.period, master_nodes.0.period_unit,
      master_nodes.0.auto_renew, streaming_core_nodes.0.charging_mode, streaming_core_nodes.0.period, streaming_core_nodes.0.period_unit,
      streaming_core_nodes.0.auto_renew, analysis_core_nodes.0.charging_mode, analysis_core_nodes.0.period, analysis_core_nodes.0.period_unit,
      analysis_core_nodes.0.auto_renew, custom_nodes.0.charging_mode, custom_nodes.0.period, custom_nodes.0.period_unit, custom_nodes.0.auto_renew
    ]
  }
}
```
