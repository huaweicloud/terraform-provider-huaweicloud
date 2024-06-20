---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_slow_log_link"
description: |-
  Use this data source to get the RDS slow log link.
---

# huaweicloud_rds_slow_log_link

Use this data source to get the RDS slow log link.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_slow_log_link" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource. If omitted, the provider-level
  region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `file_name` - (Optional, String) Specifies the name of the file to be downloaded. It is mandatory when the type of the
  instance is **SQLServer**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `file_size` - Indicates the file size in KB.

* `file_link` - Indicates the download link.

* `created_at` - Indicates the creation time.

## Timeouts

This data source provides the following timeouts configuration options:

* `read` - Default is 10 minutes.
