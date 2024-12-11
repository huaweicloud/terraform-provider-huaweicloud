---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_rds_database"
description: |-
  Manage the resource of adding RDS database to DBSS instance within HuaweiCloud.
---

# huaweicloud_dbss_rds_database

Manage the resource of adding RDS database to DBSS instance within HuaweiCloud.

-> Before adding the RDS database to the DBSS instance, the DBSS instance `status` must be **ACTIVE**.

## Example Usage

```hcl
variable "instance_id" {}
variable "rds_id" {}
variable "type" {}

resource "huaweicloud_dbss_rds_database" "test" {
  instance_id = var.instance_id
  rds_id      = var.rds_id
  type        = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the DBSS instance ID.
  Changing this parameter will create a new resource.

* `rds_id` - (Required, String, ForceNew) Specifies the RDS instance ID.
  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the RDS database type.
  The valid values are as follows:
  + **MYSQL**
  + **ORACLE**
  + **POSTGRESQL**
  + **SQLSERVER**
  + **DAMENG**
  + **TAURUS**
  + **DWS**
  + **KINGBASE**
  + **MARIADB**
  + **GAUSSDBOPENGAUSS**

  Changing this parameter will create a new resource.

* `status` - (Optional, String) Specifies the audit status of the RDS database.
  The valid values are as follows:
  + **ON**
  + **OFF**

  After an RDS database is associated with the DBSS instance, the audit status is **OFF** by default.

* `lts_audit_switch` - (Optional, Int) Specifies whether to disable LTS audit.
  The valid values are as follows:
  + `1`: Indicates disable.
  + `0`: Remain unchanged. (In this case, the value can also be an integer other than `1`).

  -> This parameter is used in the DWS database scenario. If you do not need to close it,
    there is no need to set this field.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `db_id` - The database ID.

* `name` - The database name.

* `version` - The database version.

* `charset` - The database character set.
  The value can be **GBK** or **UTF8**

* `ip` - The database IP address.

* `port` - The database port.

* `os` - The database operation system.

* `instance_name` - The database instance name.

* `audit_status` - The database running status.
  The value can be **ACTIVE**, **SHUTOFF** or **ERROR**.

* `agent_url` - The unique ID of the agent.

* `db_classification` - The classification of the database.
  The value can be **RDS** (RDS database) or **ECS** (self-built database).

* `rds_audit_switch_mismatch` - Whether the audit switch status of the RDS instance is match.
  When the database audit function is enabled and the log upload function on RDS is disabled, the value is **true**.

## Import

The resource can be imported using the related `instance_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dbss_rds_database.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response.
The missing attributes include: `lts_audit_switch`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to align
with the instance. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_dbss_rds_database" "test" {
  ...

  lifecycle {
    ignore_changes = [
      lts_audit_switch,
    ]
  }
}
```
