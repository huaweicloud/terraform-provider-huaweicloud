---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_databases"
description: |-
  Use this data source to get a list of databases under a specified DBSS instance.
---

# huaweicloud_dbss_databases

Use this data source to get a list of databases under a specified DBSS instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dbss_databases" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the audit instance ID to which the databases belong.

* `status` - (Optional, String) Specifies the status of the database.
  The value can be **ON** or **OFF**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - The list of the databases.

  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `id` - The ID of the database.

* `name` - The name of the database.

* `status` - The status of the database.

* `type` - The type of the added database.
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

* `ip` - The IP address of the database.

* `os` - The operation system of the database.

* `port` - The port of the database.

* `version` - The version of the database.

* `charset` - The character set of the database.

* `instance_name` - The name of the database instance.

* `agent_url` - The unique ID of the agent.

* `audit_status` - The running status of the database.
  The value can be **ACTIVE**, **SHUTOFF** or **ERROR**.

* `db_classification` - The classification of the database.
  The value can be **RDS** (RDS database) or **ECS** (self-built database).
