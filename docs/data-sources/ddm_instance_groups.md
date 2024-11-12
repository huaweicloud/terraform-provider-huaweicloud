---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_instance_groups"
description: |-
  Use this data source to get the list of DDM instance groups.
---

# huaweicloud_ddm_instance_groups

Use this data source to get the list of DDM instance groups.

## Example Usage

```hcl
variable "ddm_instance_id" {}

data "huaweicloud_ddm_instance_groups" "test" {
  instance_id = var.ddm_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of DDM instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `group_list` - Indicates the list of DDM instance group.
  The [group_list](#group_list_struct) structure is documented below.

<a name="group_list_struct"></a>
The `group_list` block supports:

* `id` - Indicates the group ID.

* `name` - Indicates the group name.

* `role` - Indicates the group role type, which can be read/write or read-only. The value can be:
  + **rw**: read/write group.
  + **r**: read-only group.

* `endpoint` - Indicates the connection address of the group. If load balancing is not enabled, the connection address
  string of the node in the group is returned.

* `ipv6_endpoint` - Indicates the IPv6 connection address of the group.

* `is_load_balance` - Indicates whether load balancing is enabled.

* `is_default_group` - Indicates whether the API group is the default group.

* `cpu_num_per_node` - Indicates the number of CPU cores per node.

* `mem_num_per_node` - Indicates the memory size per node, in GB.

* `architecture` - Indicates the CPU architecture. The value can be: **x86**, **Arm**.

* `node_list` - Indicates the list of node.
  The [node_list](#group_list_node_list_struct) structure is documented below.

<a name="group_list_node_list_struct"></a>
The `node_list` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `az` - Indicates the AZ to which the node belongs.
