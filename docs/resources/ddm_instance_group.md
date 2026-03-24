---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_instance_group"
description: |-
  Manages a DDM instance group resource within HuaweiCloud.
---

# huaweicloud_ddm_instance_group

Manages a DDM instance group resource within HuaweiCloud.

## Example Usage

```hcl
variable instance_id {}
variable flavor_id {}
variable subnet_id {}

resource "huaweicloud_ddm_instance_group" "test" {
  instance_id = var.instance_id
  name        = "test_group_name"
  type        = "rw"
  flavor_id   = var.flavor_id

  nodes {
    available_zone = "cn-north-4a"
    subnet_id      = var.subnet_id
  }
  nodes {
    available_zone = "cn-north-4b"
    subnet_id      = var.subnet_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of a DDM instance.

* `name` - (Required, String, NonUpdatable) Specifies the name of a DDM instance group, it can include 4 to 64 characters,
  must start with a letter and can contain only letters, digits, underscores (_), and hyphens (-).

* `type` - (Required, String, NonUpdatable) Specifies the type of the instance group. Value options:
  + **rw**: read/write group.
  + **r**: read-only group.

* `flavor_id` - (Required, String, NonUpdatable) Specifies the specification ID.

* `nodes` - (Required, List, NonUpdatable) Specifies the node information list.
  The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `available_zone` - (Required, String, NonUpdatable) Specifies the AZ where the node is located.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the subnet ID.

## Attribute Reference

* `id` - The resource ID.

* `endpoint` - Indicates the connection address of the group. If load balancing is not enabled, the connection address
  string of the node in the group is returned.

* `ipv6_endpoint` - Indicates the IPv6 connection address of the group.

* `is_load_balance` - Indicates whether load balancing is enabled.

* `is_default_group` - Indicates whether the API group is the default group.

* `cpu_num_per_node` - Indicates number of CPU cores per node.

* `mem_num_per_node` - Indicates memory size per node, in GB.

* `architecture` - Indicates the CPU architecture. The value can be **x86** and **Arm**.

* `nodes` - Indicates the node information list.
  The [nodes](#nodes_attribute) structure is documented below.

<a name="nodes_attribute"></a>
The `nodes` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The DDM instance group can be imported using the `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_ddm_instance_group.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `flavor_id` and `nodes/subnet_id`. It is
generally recommended running `terraform plan` after importing a DDM instance group. You can then decide if changes
should be applied to the DDM instance group, or the resource definition should be updated to align with the DDM instance
group. Also you can ignore changes as below.

```hcl
resource "huaweicloud_ddm_instance_group" "test" {
  ...

  lifecycle {
    ignore_changes = [
      flavor_id, nodes.0.subnet_id
    ]
  }
}
```
