---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_schema"
description: |-
  Manages a GaussDB OpenGauss schema resource within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_schema

Manages a GaussDB OpenGauss schema resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_opengauss_schema" "test" {
  instance_id = var.instance_id
  db_name     = "test_db_name"
  name        = "test_schema_name"
  owner       = "root"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB OpenGauss instance.

  Changing this parameter will create a new resource.

* `db_name` - (Required, String, ForceNew) Specifies the database name.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the schema name. The value can contain **1** to **63** characters. Only
  letters, digits, and underscores (_) are allowed. It cannot start with pg or a digit, and must be different from template
  database names and existing schema names. Template databases include **postgres**, **template0**, **template1**.
  Existing schemas include **public** and **information_schema**.

  Changing this parameter will create a new resource.

* `owner` - (Required, String, ForceNew) Specifies the owner of the schema. The value cannot be a system user and must be
  an existing database username. System users: **rdsAdmin**, **rdsMetric**, **rdsBackup**, and **rdsRepl**.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<db_name>/<name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 90 minutes.
* `delete` - Default is 90 minutes.

## Import

The GaussDB OpenGauss schema can be imported using the `instance_id`, `db_name` and `name` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_gaussdb_opengauss_schema.test <instance_id>/<db_name>/<name>
```
