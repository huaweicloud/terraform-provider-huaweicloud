---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_database"
description: |-
  Manages a GaussDB database resource within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_database

Manages a GaussDB database resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_opengauss_database" "test" {
  instance_id = var.instance_id
  name        = "test_db_name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the database name. The name can contain `1` to `63` characters.
  Only letters, digits, hyphens (-), and underscores (_) are allowed. It cannot start with **pg** or a digit, and must
  be different from template database names. Template databases include **postgres**, **template0**, **template1** and
  **templatem**.

  Changing this parameter will create a new resource.

* `character_set` - (Optional, String, ForceNew) Specifies the database character set. Defaults to **UTF8**.

  Changing this parameter will create a new resource.

* `owner` - (Optional, String, ForceNew) Specifies the Database user. Defaults to **root**. The value must be an existing
  username and must be different from system usernames. System users: **rdsAdmin**, **rdsMetric**, **rdsBackup** and
  **rdsRepl**.

  Changing this parameter will create a new resource.

* `template` - (Optional, String, ForceNew) Specifies the name of the database template. The value can be **template0**.

  Changing this parameter will create a new resource.

* `lc_collate` - (Optional, String, ForceNew) Specifies the database collation. Defaults to **C**.

  Changing this parameter will create a new resource.

* `lc_ctype` - (Optional, String, ForceNew) Specifies the database classification. Defaults to **C**.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<name>`.

* `size` - Indicates the database size.

* `compatibility_type` - Indicates the database compatibility type.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 90 minutes.
* `delete` - Default is 90 minutes.

## Import

The GaussDB database can be imported using the `instance_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_opengauss_database.test <instance_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `template` and `lc_ctype`. It is generally
recommended running `terraform plan` after importing a GaussDB database. You can then decide if changes should be applied
to the GaussDB database, or the resource definition should be updated to align with the GaussDB database. Also, you can
ignore changes as below.

```hcl
resource "huaweicloud_gaussdb_opengauss_database" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      template, lc_ctype,
    ]
  }
}
```
