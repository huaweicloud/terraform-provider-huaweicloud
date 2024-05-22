---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_databases"
description: |-
  Use this data source to get the list of GaussDB MySQL databases.
---

# huaweicloud_gaussdb_mysql_databases

Use this data source to get the list of GaussDB MySQL databases.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_mysql_databases" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of GaussDB MySQL Instance.

* `name` - (Optional, String) Specifies the database name.

* `character_set` - (Optional, String) Specifies the database character set,
  Value options: **utf8mb4**, **utf8**, **latin1**, **gbk**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - Indicates the list databases.

  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `name` - Indicates the  database name.

* `character_set` - Indicates the  database character set.

* `description` - Indicates the  database comment.

* `users` - Indicates the list of authorized database users.

  The [users](#databases_users_struct) structure is documented below.

<a name="databases_users_struct"></a>
The `users` block supports:

* `name` - Indicates the  database username.

* `host` - Indicates the  host IP address.

* `readonly` - Indicates whether the database permission is read-only.
