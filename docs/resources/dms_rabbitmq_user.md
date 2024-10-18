---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_user"
description: |-
  Manages a DMS RabbitMQ user resource within HuaweiCloud.
---

# huaweicloud_dms_rabbitmq_user

Manages a DMS RabbitMQ user resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "user_name" {}

resource "huaweicloud_dms_rabbitmq_user" "test" {
  instance_id = var.instance_id
  access_key  = var.user_name
  secret_key  = "Terraform@123"

  vhosts {
    vhost = "default"
    conf  = ".*"
    write = ".*"
    read  = ".*"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the RabbitMQ instance ID.
  Changing this creates a new resource.

* `access_key` - (Required, String, ForceNew) Specifies the user name. It starts with a letter, consists of 7 to 64
  characters, and contains only letters, digits, hyphens (-), and underscores (_).
  Changing this creates a new resource.

* `secret_key` - (Required, String) Specifies the user password. It consists of 8 to 32 characters.
  Contain at least three of the following character types:
  + Uppercase letters
  + Lowercase letters
  + Digits
  + Special characters `~!@#$%^&*()-_=+\|[{}];:'",<.>/?

  It cannot be the user name or the user name spelled backwards.

* `vhosts` - (Required, List) Specifies the virtual hosts to be granted permissions for.
  The [vhosts](#block--vhosts) structure is documented below.

<a name="block--vhosts"></a>
The `vhosts` block supports:

* `vhost` - (Required, String) Specifies the name of the virtual host to be granted permissions for.

* `conf` - (Required, String) Specifies the granting resource permissions using regular expressions.

* `read` - (Required, String) Specifies the granting resource read permissions using regular expressions.

* `write` - (Required, String) Specifies the granting resource write permissions using regular expressions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The RabbitMQ user can be imported using `instance_id` and `access_key` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dms_rabbitmq_user.test <instance_id>/<access_key>
```
