---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_accounts"
description: ""
---

# huaweicloud_rds_mysql_accounts

Use this data source to get the list of RDS MySQL accounts.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_rds_mysql_accounts" "test" {
  instance_id   = var.instance_id
  name          = "test"
  host          = "10.10.10.10"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `name` - (Optional, String) Specifies the username of the DB account.

* `host` - (Optional, String) Specifies the IP address that is allowed to access your DB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - Indicates the list of users.
  The [users](#RDS_mysql_users) structure is documented below.

<a name="RDS_mysql_users"></a>
The `users` block supports:

* `name` - Indicates the username of the DB account.

* `hosts` - Indicates the IP addresses that are allowed to access your DB instance.

* `description` - Indicates remarks of the database account.
