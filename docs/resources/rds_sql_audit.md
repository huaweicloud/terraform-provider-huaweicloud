---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sql_audit"
description: ""
---

# huaweicloud_rds_sql_audit

Manages RDS SQL audit resource within HuaweiCloud.

-> **NOTE:** Only MySQL and PostgreSQL engines are supported.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_sql_audit" "test" {
  instance_id = var.instance_id
  keep_days   = 5
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.

* `keep_days` - (Required, Int) Specifies the number of days for storing audit logs. Value ranges from `1` to `732`.

* `audit_types` - (Optional, List) Specifies the list of audit types. Value options: **CREATE_USER**, **DROP_USER**,
  **RENAME_USER**, **GRANT**, **REVOKE**, **CREATE**, **ALTER**, **DROP**, **RENAME**, **TRUNCATE**, **INSERT**,
  **DELETE**, **UPDATE**, **REPLACE**, **SELECT**, **BEGIN/COMMIT/ROLLBACK**, **PREPARED_STATEMENT**.
  It is not supported for PostgreSQL.

* `reserve_auditlogs` - (Optional, Bool) Specifies whether the historical audit logs will be reserved for some time
  when SQL audit is disabled. It is valid only when SQL audit is disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS SQL audit can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rds_sql_audit.test <id>
```
