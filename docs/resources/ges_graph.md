---
subcategory: "Graph Engine Service (GES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ges_graph"
description: ""
---

# huaweicloud_ges_graph

Manages a GES graph resource within HuaweiCloud.  

## Example Usage

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
  
resource "huaweicloud_ges_graph" "test" {
  name                  = "demo"
  graph_size_type_index = "1"
  cpu_arch              = "x86_64"
  vpc_id                = var.vpc_id
  subnet_id             = var.subnet_id
  security_group_id     = var.secgroup_id
  crypt_algorithm       = "generalCipher"
  enable_https          = false

  tags = {
    key = "val"
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) The graph name.  
  The name must start with a letter and contains 4 to 50 characters consisting of letters,
  digits, hyphens (-), and underscores (_). It cannot contain special characters.

  Changing this parameter will create a new resource.

* `graph_size_type_index` - (Required, String) Graph size type index.  
  Value options are as follows:
   + **0**: indicates 10 thousand edges.
   + **1**: indicates 1 million edges.
   + **2**: indicates 10 million edges.
   + **3**: indicates 100 million edges.
   + **4**: indicates 1 billion edges.
   + **5**: indicates 10 billion edges.
   + **6**: indicates the database edition.
   + **401**: indicates 1 billion enhanced edges.

* `vpc_id` - (Required, String, ForceNew) The VPC ID.

  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) The subnet ID.

  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String, ForceNew) The security group ID.

  Changing this parameter will create a new resource.

* `crypt_algorithm` - (Required, String, ForceNew) Graph instance cryptography algorithm.  
  Value options are as follows:
    + **generalCipher**: Chinese cryptographic algorithm.
    + **SMcompatible**: Commercial cryptography algorithm (compatible with international ones).

  Changing this parameter will create a new resource.

* `enable_https` - (Required, Bool, ForceNew) Whether to enable the security mode. This mode may damage GES performance
  greatly.

  Changing this parameter will create a new resource.

* `cpu_arch` - (Optional, String, ForceNew) Graph instance's CPU architecture type.  
 The value can be **x86_64** or **aarch64**. The default value is **x86_64**.

  Changing this parameter will create a new resource.

* `public_ip` - (Optional, List, ForceNew) The information about public IP.  
  If the parameter is not specified, public connection is not used by default.

  Changing this parameter will create a new resource.

  The [PublicIp](#GesGraph_PublicIp) structure is documented below.

* `enable_multi_az` - (Optional, Bool, ForceNew) Whether the created graph supports the cross-AZ mode.
  The default value is false.  
  If the value is true, the system will create the ECSs in the graph in two AZs.

  Changing this parameter will create a new resource.

* `encryption` - (Optional, List, ForceNew) The configuration of data encryption.

  Changing this parameter will create a new resource.

  The [Encryption](#GesGraph_Encryption) structure is documented below.

* `lts_operation_trace` - (Optional, List, ForceNew) The configuration of audit logs.

  Changing this parameter will create a new resource.

  The [LtsOperationTrace](#GesGraph_LtsOperationTrace) structure is documented below.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project ID.

  Changing this parameter will create a new resource.

* `tags` - (Optional, Map, ForceNew) The key/value pairs to associate with the graph.

  Changing this parameter will create a new resource.

* `enable_rbac` - (Optional, Bool, ForceNew) Whether to enable granular permission control for the created graph.  
  The default value is false. If this parameter is set to true, no user has the permission to access the graph.
  To access the graph, you need to call the granular permission control API of the service plane
   to set the required permissions.

  Changing this parameter will create a new resource.

* `enable_full_text_index` - (Optional, Bool, ForceNew) Whether to enable full-text index control for the created graph.
  The default value is false. If this parameter is set to true, full-text indexes are available
  for 1-billion-edge-pro graphs, and a Cloud Search Service (CSS) cluster will
  be created when you create a graph.

  -> If you enable full-text indexes: If the CSS has been deployed, the system automatically creates a
   CSS cluster during the creation of the graph instance, which will take a long time.
   If the CSS is not deployed, the graph creation will fail.

  Changing this parameter will create a new resource.

* `enable_hyg` - (Optional, Bool, ForceNew) Whether to enable HyG for the graph.
  This parameter is available for database edition graphs only.

  Changing this parameter will create a new resource.

* `product_type` - (Optional, String, ForceNew) Graph product type.  
  Value options are as follows:
    + **InMemory**: memory edition.
    + **Persistence**: database edition.
  
  If **graph_size_type_index** is 6, the value must be **Persistence**.

  Changing this parameter will create a new resource.

* `vertex_id_type` - (Optional, List, ForceNew) The configuration of vertex ID.
  This parameter is mandatory only for database edition graphs.

  Changing this parameter will create a new resource.

  The [vertexIdType](#GesGraph_vertexIdType) structure is documented below.

* `replication` - (Optional, Int) Number of replicas.

* `keep_backup` - (Optional, Bool, ForceNew) Whether to retain the backups of a graph after it is deleted.

<a name="GesGraph_PublicIp"></a>
The `PublicIp` block supports:

* `public_bind_type` - (Optional, String) The bind type of public IP.  
  The valid value are **auto_assign**, and **bind_existing**.

* `eip_id` - (Optional, String) The EIP ID.  

<a name="GesGraph_Encryption"></a>
The `Encryption` block supports:

* `enable` - (Optional, Bool) Whether to enable data encryption. The value can be true or false.
  The default value is false.  

* `master_key_id` - (Optional, String) ID of the customer master key created by DEW in the project corresponding
  to the graph creation.  

<a name="GesGraph_LtsOperationTrace"></a>
The `LtsOperationTrace` block supports:

* `enable_audit` - (Optional, Bool) Whether to enable graph audit. The default value is false.  

* `audit_log_group_name` - (Optional, String) LTS log group name.  

<a name="GesGraph_vertexIdType"></a>
The `vertexIdType` block supports:

* `id_type` - (Optional, String) Vertex ID type.  
  Value options are as follows:
    + **fixedLengthString**: Vertex IDs are used for internal storage and compute.
        Specify the length limit. If the IDs are too long, the query performance can be reduced.
        Specify the length limit based on your dataset vertex IDs.
    + **hash**: Vertex IDs are converted into hash code for storage and compute.
        There is no limit on the ID length. However, there is an extremely low probability, approximately 10^(-43),
        that the vertex IDs will conflict. If you cannot determine the maximum length of a vertex ID,
        set this parameter to Hash.

* `id_length` - (Optional, Int) The length of ID.  
  This parameter is mandatory if **id_type** is **fixedLengthString**. The value ranges from 1 to 128.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `az_code` - AZ code

* `status` - Status of a graph.  
  The value can be one of the following:
    + **100**: Indicates that a graph is being prepared.
    + **200**: indicates that a graph is running.
    + **201**: indicates that a graph is upgrading.
    + **202**: indicates that a graph is being imported.
    + **203**: indicates that a graph is being rolled back.
    + **204**: indicates that a graph is being exported.
    + **205**: indicates that a graph is being cleared.
    + **206**: indicates that the system is preparing for resize.
    + **207**: indicates that the resize is in progress.
    + **208**: Indicates that the resize is being rolled back.
    + **210**: Preparing for expansion
    + **211**: Expanding
    + **300**: indicates that a graph is faulty.
    + **303**: indicates that a graph fails to be created.
    + **400**: indicates that a graph is deleted.
    + **800**: indicates that a graph is frozen.
    + **900**: indicates that a graph is stopped.
    + **901**: indicates that a graph is being stopped.
    + **920**: indicates that a graph is being started.

* `private_ip` - Floating IP address of a graph instance.

* `traffic_ip_list` - Physical addresses of a graph instance for access from private networks.  
  To prevent service interruption caused by floating IP address switchover,
  poll the physical IP addresses to access the graph instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

The ges graph can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ges_graph.test a0be840b-b223-48da-8b34-b8fee1b2e0ca
```
