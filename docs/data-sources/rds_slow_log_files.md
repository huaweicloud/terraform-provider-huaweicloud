---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_slow_log_files"
description: |-
  Use this data source to get the list of RDS slow log files.
---

# huaweicloud_rds_slow_log_files

Use this data source to get the list of RDS slow log files.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_slow_log_files" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `files` - Indicates the list of slow log files.

  The [files](#files_struct) structure is documented below.

<a name="files_struct"></a>
The `files` block supports:

* `file_name` - Indicates the file name.

* `file_size` - Indicates the file size in bytes.
