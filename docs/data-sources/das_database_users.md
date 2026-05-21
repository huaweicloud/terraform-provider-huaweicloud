---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_database_users"
description: |-
  Use this data source to query DAS Database users within HuaweiCloud.
---

# huaweicloud_das_database_users

Use this data source to query DAS Database users within HuaweiCloud.

## Example Usage

### Query all database users under a specified instance

```hcl
variable "instance_id" {}

data "huaweicloud_das_database_users" "test" {
  instance_id = var.instance_id
}
```

### Query a specific database user under a specified instance by user name

```hcl
variable "instance_id" {}
variable "user_name" {}

data "huaweicloud_das_database_users" "test" {
  instance_id = var.instance_id
  user_name   = var.user_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the database users are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the instance to which the database user belongs.

* `user_id` - (Optional, String) Specifies the ID of the database user.

* `user_name` - (Optional, String) Specifies the name of the database user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `users` - The list of users that matched filter parameters.  
  The [users](#das_database_users) structure is documented below.

<a name="das_database_users"></a>
The `users` block supports:

* `id` - The ID of the database user, in UUID format.

* `name` - The name of the database user.
