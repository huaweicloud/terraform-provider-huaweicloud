---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_pt_applicable_instances"
description: |-
  Use this data source to get the list of DDS parameter template applicable instances.
---

# huaweicloud_dds_pt_applicable_instances

Use this data source to get the list of DDS parameter template applicable instances.

## Example Usage

```hcl
variable "configuration_id"  {}

data "huaweicloud_dds_pt_applicable_instances" "test" {
  configuration_id = var.configuration_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `configuration_id` - (Required, String) Specifies the ID of the parameter template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the applicable instances.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `entities` - Indicates the list of node group information or node information.
  + For a replica set instance, it is not returned and the parameter template is directly applied to the corresponding
  DB instance.
  + For the shard or config group of a cluster instance, information about the node group of the corresponding instance
  is returned. For the mongos group of a cluster instance, information about the node of the corresponding instance is
  returned.

  The [entities](#instances_entities_struct) structure is documented below.

<a name="instances_entities_struct"></a>
The `entities` block supports:

* `entity_id` - Indicates the group ID or node ID.

* `entity_name` - Indicates the group name or node name.
