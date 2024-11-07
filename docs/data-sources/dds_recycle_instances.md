---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_recycle_instances"
description: |-
  Use this data source to get the list of DDS recycle instances.
---

# huaweicloud_dds_recycle_instances

Use this data source to get the list of DDS recycle instances.

## Example Usage

```hcl
data "huaweicloud_dds_recycle_instances" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the instances.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - Indicates the instance ID.

* `name` - Indicates the instance name.

* `mode` - Indicates the instance mode.

* `backup_id` - Indicates the backup ID.

* `datastore` - Indicates the database information.

  The [datastore](#instances_datastore_struct) structure is documented below.

* `charging_mode` - Indicates the charging mode.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `created_at` - Indicates the creation time.

* `deleted_at` - Indicates the deletion time.

* `retained_until` - Indicates the retention end time.

<a name="instances_datastore_struct"></a>
The `datastore` block supports:

* `version` - Indicates the database version.

* `type` - Indicates the database type.
