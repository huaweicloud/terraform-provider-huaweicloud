---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_auto_launch_group_instances"
description: |-
  Use this data source to get the list of instances created by the auto launch group.
---

# huaweicloud_compute_auto_launch_group_instances

Use this data source to get the list of instances created by the auto launch group.

## Example Usage

```hcl
variable "auto_launch_group_id" {}

data "huaweicloud_compute_auto_launch_group_instances" "test" {
  auto_launch_group_id = var.auto_launch_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `auto_launch_group_id` - (Required, String) Specifies the ID of the auto launch group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of the instance created by the auto launch group.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - Indicates the ID of the instance.

* `name` - Indicates name ID of the instance.

* `flavor_id` - Indicates flavor ID of the instance.

* `availability_zone_id` - Indicates the ID of the AZ.

* `status` - Indicates state of the instance.

* `sell_mode` - Indicates the sales model of the instance.
  The value can be: **spot**, **onDemand**.
