---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_database_user"
description: |-
  Manages DAS Database user resource within HuaweiCloud.
---

# huaweicloud_das_database_user

Manages DAS Database user resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "user_name" {}
variable "user_password" {}

resource "huaweicloud_das_database_user" "test" {
  instance_id = var.instance_id
  name        = var.user_name
  password    = var.user_password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the database user is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the instance to which the database user belongs.

* `name` - (Required, String) Specifies the name of the database user.

* `password` - (Required, String) Specifies the password of the database user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

## Import

The database user can be imported using `instance_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_das_database_user.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the imported state. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_das_database_user" "test" {
  ...

  lifecycle {
    ignore_changes = [
      password,
    ]
  }
}
```
