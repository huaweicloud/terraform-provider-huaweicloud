---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_account_roles"
description: |-
  Manages an RDS PostgreSQL account roles resource within HuaweiCloud.
---

# huaweicloud_rds_pg_account_roles

Manages an RDS PostgreSQL account roles resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_pg_account_roles" "test" {
  instance_id = var.instance_id
  user        = "test_user"
  roles       = ["test111", "test222"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS PostgreSQL instance.

* `user` - (Required, String, NonUpdatable) Specifies the username of the account.

* `roles` - (Required, List) Specifies the list of roles.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of account which is formatted `<instance_id>/<name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS PostgreSQL account roles can be imported using the `instance_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rds_pg_account_roles.test <instance_id>/<name>
```
