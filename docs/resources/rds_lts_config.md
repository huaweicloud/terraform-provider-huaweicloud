---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_lts_config"
description: |-
  Manages a RDS LTS config resource within HuaweiCloud.
---

# huaweicloud_rds_lts_config

Manages a RDS LTS config resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "lts_group_id" {}
variable "lts_stream_id" {}

resource "huaweicloud_rds_lts_config" "test" {
  instance_id   = var.instance_id
  engine        = "mysql"
  log_type      = "error_log"
  lts_group_id  = var.lts_group_id
  lts_stream_id = var.lts_stream_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `engine` - (Required, String, ForceNew) Specifies the engine of the RDS instance.
  Value options: **mysql**, **postgresql**, **sqlserver**. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RDS instance.
  Changing this creates a new resource.

* `log_type` - (Required, String, ForceNew) Specifies the type of the LTS log config.
  Value options: **error_log**, **slow_log**, **audit_log**. Changing this creates a new resource.

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

The RDS LTS log can be imported using `instance_id` and `log_type` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rds_lts_config.test <instance_id>/<log_type>
```
