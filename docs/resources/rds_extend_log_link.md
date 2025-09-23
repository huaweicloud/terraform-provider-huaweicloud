---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_extend_log_link"
description: |-
  Manages an RDS extend log link resource within HuaweiCloud.
---

# huaweicloud_rds_extend_log_link

Manages an RDS extend log link resource within HuaweiCloud.  

## Example Usage

```hcl
variable "instance_id" {}
variable "file_name" {}

resource "huaweicloud_rds_extend_log_link" "test" {
  instance_id = var.instance_id
  file_name   = var.file_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.

* `file_name` - (Required, String, NonUpdatable) Specifies the name of the file to be downloaded.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The format is  `<instance_id>/<file_name>`.

* `file_size` - Indicates the file size in KB.

* `file_link` - Indicates the download link.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
