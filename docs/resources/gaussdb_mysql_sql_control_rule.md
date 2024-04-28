---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_sql_control_rule"
description: ""
---

# huaweicloud_gaussdb_mysql_sql_control_rule

Manages a GaussDB MySQL SQL concurrency control rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_gaussdb_mysql_sql_control_rule" "test" {
  instance_id     = var.instance_id
  node_id         = var.node_id
  sql_type        = "SELECT"
  pattern         = "select~from~t1"
  max_concurrency = 20
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance.

  Changing this parameter will create a new resource.

* `node_id` - (Required, String, ForceNew) Specifies ID of the GaussDB redis node.

  Changing this parameter will create a new resource.

* `sql_type` - (Required, String, ForceNew) Specifies SQL statement type.
  Value options: **SELECT**, **UPDATE**, **DELETE**.

  Changing this parameter will create a new resource.

* `pattern` - (Required, String, ForceNew) Specifies the concurrency control rule of SQL statements. A rule can consist
  of up to 128 keywords. The keywords are separated by tildes (~), for example, select~from~t1. The rule cannot contain
  backslashes (\), commas (,), or double tildes (~~). It cannot end with tildes (~).

  Changing this parameter will create a new resource.

* `max_concurrency` - (Required, Int) Specifies the maximum number of concurrent SQL statements.
  Value: a non-negative integer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<node_id>/<sql_type>/<pattern>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The GaussDB MySQL SQL concurrency control rule can be imported using the `instance_id`,`node_id`,`sql_type` and
`pattern` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_sql_control_rule.test <instance_id>/<node_id>/<sql_type>/<pattern>
```
