---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_lts_log"
description: |-
  Manages a GaussDB MySQL LTS log resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_lts_log

Manages a GaussDB MySQL LTS log resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "lts_group_id" {}
variable "lts_stream_id" {}

resource "huaweicloud_gaussdb_mysql_lts_log" "test" {
  instance_id   = var.instance_id
  log_type      = "error_log"
  lts_group_id  = var.lts_group_id
  lts_stream_id = var.lts_stream_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance.
  Changing this creates a new resource.

* `log_type` - (Required, String, ForceNew) Specifies the type of the LTS log.
  Value options: **error_log**, **slow_log**. Changing this creates a new resource.

* `lts_group_id` - (Required, String) Specifies the ID of the LTS log group.

* `lts_stream_id` - (Required, String) Specifies the ID of the LTS log stream.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID in format of `<instance_id>/<log_type>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The GaussDB MySQL LTS log can be imported using `instance_id` and `log_type` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_lts_log.test <instance_id>/<log_type>
```
