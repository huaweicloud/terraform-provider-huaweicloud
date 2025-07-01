---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_binlog"
description: ""
---

# huaweicloud_rds_mysql_binlog

Manages RDS MySQL binlog resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_mysql_binlog" "test" {
  instance_id            = var.instance_id
  binlog_retention_hours = 6
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the RDS binlog resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the RDS MySQL instance ID.

* `binlog_retention_hours` - (Required, Int) Specifies the binlog retention period. Value range: `1` to `168` (7x24).

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID .

## Import

RDS MySQL binlog can be imported using the `instance id`, e.g.

```bash
$ terraform import huaweicloud_rds_mysql_binlog.test <instance_id>
```
