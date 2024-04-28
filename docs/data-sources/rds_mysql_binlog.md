---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_binlog"
description: ""
---

# huaweicloud_rds_mysql_binlog

Use this data source to get the binlog retention hours of RDS MySQL.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_rds_mysql_binlog" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS MySQL instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `binlog_retention_hours` - The binlog retention period. Value range: 0 to 168 (7x24).
