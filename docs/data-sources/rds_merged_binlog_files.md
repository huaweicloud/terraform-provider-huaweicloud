---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_merged_binlog_files"
description: |-
  Use this data source to query the merged binlog files within HuaweiCloud.
---

# huaweicloud_rds_merged_binlog_files

Use this data source to query the merged binlog files within HuaweiCloud.

-> This interface only supports **MySQL** engine.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_merged_binlog_files" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the pack log infos.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pack_log_infos` - The list of Binlog merged download files.
  The [pack_log_infos](#pack_log_infos_attr) structure is documented below.

<a name="pack_log_infos_attr"></a>
The `pack_log_infos` block supports:

* `id` - The unique ID of the file.

* `instance_id` - The ID of the RDS instance.

* `size` - The file size.

* `size_unit` - The unit of the file size.

* `status` - The status of the file.

* `query_start_time` - The start timestamp of the merge time range.

* `query_end_time` - The end timestamp of the merge time range.

* `file_name` - The name of the file.
