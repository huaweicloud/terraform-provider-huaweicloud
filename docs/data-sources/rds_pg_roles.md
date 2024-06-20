---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_roles"
description: |-
  Use this data source to get the list of RDS instance roles.
---

# huaweicloud_rds_pg_roles

Use this data source to get the list of RDS instance roles.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_instance_roles" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of Instance.

* `account` - (Optional, String) Specifies the username of the account.
  The list of authorized roles for this account when there is a value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `roles` - Indicates the list of roles.
