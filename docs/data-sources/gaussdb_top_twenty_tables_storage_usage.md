---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_top_twenty_tables_storage_usage"
description: |-
  Use this data source to query the storage usage of the top 20 tables of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_top_twenty_tables_storage_usage

Use this data source to query the storage usage of the top 20 tables of a GaussDB instance within HuaweiCloud.

## Example Usage

### Query top 20 tables by storage usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_top_twenty_tables_storage_usage" "test" {
  instance_id = var.instance_id
}
```

### Query with job_id and node_id

```hcl
variable "instance_id" {}
variable "job_id" {}
variable "node_id" {}

data "huaweicloud_gaussdb_top_twenty_tables_storage_usage" "test" {
  instance_id = var.instance_id
  job_id      = var.job_id
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the top table storage usage.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `job_id` - (Optional, String) Specifies the workflow ID, obtained from the first call without any task parameters.

* `node_id` - (Optional, String) Specifies the workflow execution node ID, obtained from the first call without any task
   parameters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `state` - The job status, which is returned when `job_id` in the request is not empty.

* `table_volumes` - The list of top table storage usage information.
  The [table_volumes](#table_volumes_struct) structure is documented below.

<a name="table_volumes_struct"></a>
The `table_volumes` block supports:

* `id` - The ID of the table.

* `table_name` - The name of the table.

* `table_owner` - The owner of the table.

* `database_name` - The name of the database.

* `schema_name` - The name of the schema.

* `is_part_type` - Whether the table has partition table properties.

* `is_hash_cluster_key` - Whether it contains hash partition column information.

* `tuples` - The number of tuples in the table.

* `create_time` - The creation time of the table.

* `update_time` - The update time of the table.

* `average_size` - The average size of the table.

* `max_ratio` - The maximum ratio.

* `min_ratio` - The minimum ratio.

* `table_size` - The size of the table (e.g., "24 kB").

* `skew_size` - The skew size of the table.

* `skew_ratio` - The skew ratio of the table.

* `skew_stddev` - The skew standard deviation of the table.
