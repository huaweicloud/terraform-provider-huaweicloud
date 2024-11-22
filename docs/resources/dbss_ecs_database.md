---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_ecs_database"
description: |-
  Manage the resource of adding self built database to DBSS instance within HuaweiCloud.
---

# huaweicloud_dbss_ecs_database

Manage the resource of adding self built database to DBSS instance within HuaweiCloud.

-> Before adding the self built database to the DBSS instance, the DBSS instance `status` must be **ACTIVE**.

## Example Usage

```hcl
variable "instance_id" {}
variable "name" {}
variable "type" {}
variable "version" {}
variable "ip" {}
variable "port" {}
variable "os" {}

resource "huaweicloud_dbss_ecs_database" "test" {
  instance_id = var.instance_id
  name        = var.name
  type        = var.type
  version     = var.version
  ip          = var.ip
  port        = var.port
  os          = var.os
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the DBSS instance ID.
  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the self built database name.
  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the self built database type.
  The valid values are as follows:
  + **MYSQL**
  + **ORACLE**
  + **POSTGRESQL**
  + **SQLSERVER**
  + **DAMENG**
  + **TAURUS**
  + **DWS**
  + **KINGBASE**
  + **GAUSSDBOPENGAUSS**
  + **GREENPLUM**
  + **HIGHGO**
  + **SHENTONG**
  + **GBASE8A**
  + **GBASE8S**
  + **GBASEXDM**
  + **MONGODB**
  + **DDS**

  Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) Specifies the self built database version.
  Changing this parameter will create a new resource.

* `ip` - (Required, String, ForceNew) Specifies the self built database IP address.
  Changing this parameter will create a new resource.

* `port` - (Required, String, ForceNew) Specifies the self built database port.
  Changing this parameter will create a new resource.

* `os` - (Required, String, ForceNew) Specifies the self built database operation system.
  The valid values are as follows:
  + **LINUX64**
  + **WINDOWS64**
  + **UNIX**

  Changing this parameter will create a new resource.

* `charset` - (Optional, String, ForceNew) Specifies the self built database character set.
  The value can be **GBK** or **UTF8**. Defaults to **UTF8**

  Changing this parameter will create a new resource.

* `instance_name` - (Optional, String, ForceNew) Specifies the self built database instance name.
  Changing this parameter will create a new resource.

* `status` - (Optional, String) Specifies the audit status of the self built database.
  The valid values are as follows:
  + **ON**
  + **OFF**

  After a self built database is associated with the DBSS instance, the audit status is **OFF** by default.

* `lts_audit_switch` - (Optional, Int) Specifies whether to disable LTS audit.
  The valid values are as follows:
  + `1`: Indicates disable.
  + `0`: Remain unchanged. (In this case, the value can also be an integer other than `1`).

  -> This parameter is used in the DWS database scenario. If you do not need to close it,
    there is no need to set this field.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `audit_status` - The database running status.
  The value can be **ACTIVE**, **SHUTOFF** or **ERROR**.

* `agent_url` - The unique ID of the agent.

* `db_classification` - The classification of the database.

## Import

The resource can be imported using the related `instance_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dbss_ecs_database.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response.
The missing attributes include: `lts_audit_switch`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to align
with the instance. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_dbss_ecs_database" "test" {
  ...

  lifecycle {
    ignore_changes = [
      lts_audit_switch,
    ]
  }
}
```
