---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sqlserver_accounts"
description: ""
---

# huaweicloud_rds_sqlserver_accounts

Use this data source to get the list of RDS SQLServer accounts.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_rds_sqlserver_accounts" "test" {
  instance_id = var.instance_id
  user_name   = "test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS SQLServer instance.

* `user_name` - (Optional, String) Specifies the username of the database account.

* `state` - (Optional, String) Specifies the database user status. Its value can be any of the following:
  + **unavailable**: The database user is unavailable.
  + **available**: The database user is available.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - Indicates the list of users.
  The [users](#RDS_sqlserver_users) structure is documented below.

<a name="RDS_sqlserver_users"></a>
The `users` block supports:

* `name` - Indicates the username of the database account.

* `state` - Indicates the database user status.
