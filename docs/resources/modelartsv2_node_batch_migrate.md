---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_node_batch_migrate"
description: |-
  Use this resource to batch migrate ModelArts nodes within HuaweiCloud.
---

# huaweicloud_modelartsv2_node_batch_migrate

Use this resource to batch migrate ModelArts nodes within HuaweiCloud.

-> This resource is a one-time action resource for batch migrating the ModelArts nodes. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "source_pool_id" {}
variable "source_cluster_id" {}
variable "target_pool_id" {}
variable "target_cluster_id" {}
variable "node_names" {
  type = list(string)
}

resource "huaweicloud_modelartsv2_node_batch_migrate" "test" {
  source_pool_id    = var.source_pool_id
  source_cluster_id = var.source_cluster_id
  target_pool_id    = var.target_pool_id
  target_cluster_id = var.target_cluster_id
  node_names        = var.node_names
          
  resource_spec {
    flavor = "modelarts.vm.cpu.16u64g.d"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource pool nodes are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `source_pool_id` - (Required, String, NonUpdatable) Specifies the ID of the source resource pool to which the
  resource nodes belong.

* `source_cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the source cluster from which
  nodes are migrated.

* `target_cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the target cluster to which
  nodes are migrated.

* `node_names` - (Optional, List, NonUpdatable) Specifies the name list of nodes to be migrated.

* `target_pool_id` - (Optional, String, NonUpdatable) Specifies the ID of the target resource pool.

* `resource_spec` - (Optional, List, NonUpdatable) Specifies the configuration information of the migrated nodes
  in the target resource pool. This parameter is **Required** for cross-pool migration.  
  The [resource_spec](#node_batch_migrate_resource_spec) structure is documented below.

<a name="node_batch_migrate_resource_spec"></a>
The `resource_spec` block supports:

* `flavor` - (Required, String) Specifies the resource specification name.  
  This parameter is **Required** for cross-pool migration.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
