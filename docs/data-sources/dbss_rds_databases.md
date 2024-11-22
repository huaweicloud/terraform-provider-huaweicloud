---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_rds_databases"
description: |-
  Use this data source to get the list of RDS databases.
---

# huaweicloud_dbss_rds_databases

Use this data source to get the list of RDS databases.

## Example Usage

```hcl
variable "type" {}

data "huaweicloud_dbss_rds_databases" "test" {
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the RDS database type.
  The valid values are as follows:
  + **MYSQL**
  + **POSTGRESQL**
  + **SQLSERVER**
  + **TAURUS**
  + **DWS**
  + **MARIADB**
  + **GAUSSDBOPENGAUSS**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - The RDS database list.

  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `id` - The RDS instance ID.

* `name` - The RDS database name.

* `type` - The RDS database type.

* `status` - The RDS instance status.
  The valid values are as follows:
  + **BUILD**: The instance is being created.
  + **ACTIVE**: The instance is normal.
  + **FAILED**: The instance is abnormal.
  + **FROZEN**: The instance is frozen.
  + **MODIFYING**: The instance is being scaled out.
  + **REBOOTING**: The instance is being restarted.
  + **RESTORING**: The instance is being restored.
  + **MODIFYING INSTANCE TYPE**: The instance is changing to the active/standby mode.
  + **SWITCHOVER**: The instance is performing an active/standby switchover.
  + **MIGRATING**: The instance is being migrated.
  + **BACKING UP**: The instance is being backed up.
  + **MODIFYING DATABASE PORT**: The database port of the instance is being changed.
  + **STORAGE FULL**: The instance disk is full.

* `version` - The RDS database version.

* `ip` - The RDS database IP address.

* `port` - The RDS database port.

* `is_supported` - Whether agent-free audit is supported.

* `instance_name` - The RDS instance name.

* `enterprise_project_id` - The enterprise project ID to which the RDS instance belongs.
