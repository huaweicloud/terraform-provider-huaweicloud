---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_user_password_reset"
description: |-
  Use this resource to reset the password of the specified user within HuaweiCloud.
---

# huaweicloud_dms_kafka_user_password_reset

Use this resource to reset the password of the specified user within HuaweiCloud.

-> This resource is only a one-time action resource for resetting the password of the specified user. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

```hcl  
variable "instance_id" {}
variable "user_name" {}
variable "new_password" {}

resource "huaweicloud_dms_kafka_user_password_reset" "test" {
  instance_id  = var.instance_id
  user_name    = var.user_name
  new_password = var.new_password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the user whose password is
  to be reset is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Kafka instance to which the user belongs.

* `user_name` - (Required, String, NonUpdatable) Specifies the name of the user to reset the password.

* `new_password` - (Required, String, NonUpdatable) Specifies the new password of the user.  
  The password has the following restrictions:
  + The password must contain `8` to `32` characters.
  + The password must contain at least three of the following: letters (case-sensitive), digits, and special
    characters (`~!@#$%^&*()-_=+|[{}]:'",<.>/?) and spaces, and cannot start with a hyphen (-).
  + The password cannot be the same as the username or the username in reverse.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
