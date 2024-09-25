---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_recycling_instances"
description: |-
  Use this data source to get the list of instances in the recycle bin.
---

# huaweicloud_gaussdb_mysql_recycling_instances

Use this data source to get the list of instances in the recycle bin.

## Example Usage

```hcl
data "huaweicloud_gaussdb_mysql_recycling_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the instances in the recycle bin.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - Indicates the instance ID.

* `name` - Indicates the instance name.

* `ha_mode` - Indicates the instance type.

* `engine_name` - Indicates the engine name.

* `engine_version` - Indicates the engine version.

* `pay_model` - Indicates the billing mode.

* `create_at` - Indicates the creation time.

* `deleted_at` - Indicates the deletion time.

* `volume_type` - Indicates the storage type.

* `volume_size` - Indicates the storage space.

* `data_vip` - Indicates the virtual IP address of the data plane.

* `data_vip_ipv6` - Indicates the IPv6 address of the data plane.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `enterprise_project_name` - Indicates the enterprise project name.

* `backup_level` - Indicates the backup level.

* `recycle_backup_id` - Indicates the backup ID.

* `recycle_status` - Indicates the recycling status.
