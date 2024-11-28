---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_snapshots"
description: |-
  Use this data source to get the list of CSS snapshots.
---

# huaweicloud_css_snapshots

Use this data source to get the list of CSS snapshots.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_css_snapshots" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backups` - The snapshot list.

  The [backups](#backups_struct) structure is documented below.

<a name="backups_struct"></a>
The `backups` block supports:

* `id` - The snapshot ID.

* `backup_expected_start_time` - The snapshot start time.

* `version` - The snapshot version.

* `restore_status` - The snapshot restoration status.

* `start_time` - The snapshot start time.

* `description` - The snapshot description.

* `backup_method` - The snapshot creation mode.

* `backup_type` - The snapshot creation type.
  The options are as follows:
  + **0**: Automatic creation.
  + **1**: Manual creation.

* `end_time` - The snapshot end time.

* `datastore` - The datastore of the cluster snapshot.

  The [datastore](#backups_datastore_struct) structure is documented below.

* `cluster_id` - The cluster ID.

* `updated_at` - The time when the snapshot was updated.

* `name` - The snapshot name.

* `status` - The snapshot status.

* `backup_keep_day` - The snapshot retention period.

* `backup_period` - The time when a snapshot is created every day.

* `indices` - The index of the back up.

* `total_shards` - The total number of shards of the back up index.

* `created_at` - The snapshot creation time.

* `cluster_name` - The cluster name.

* `failed_shards` - The number of shards that fail to be backed up.

* `bucket_name` - The name of the bucket that stores snapshot data.

<a name="backups_datastore_struct"></a>
The `datastore` block supports:

* `type` - The engine type.

* `version` - The elastic search engine version.
