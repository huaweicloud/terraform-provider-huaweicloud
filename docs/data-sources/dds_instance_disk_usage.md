---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance_disk_usage"
description: |-
  Use this data source to get the instance disk information.
---

# huaweicloud_dds_instance_disk_usage

Use this data source to get the instance disk information.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dds_instance_disk_usage" "test"{
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `volumes` - The instance disks information.

  The [volumes](#volumes_struct) structure is documented below.

<a name="volumes_struct"></a>
The `volumes` block supports:

* `entity_id` - Indicates the instance ID or group ID or node ID.

* `entity_name` - Indicates the instance name or group name or node name.

* `group_type` - Indicates the group type.
  + **mongos**
  + **shard**
  + **config**
  + **replica**
  + **single**
  + **readonly**

* `used` - Indicates the instance disk capacity used, in GB.

* `size` - Indicates the instance disk capacity, in GB.
