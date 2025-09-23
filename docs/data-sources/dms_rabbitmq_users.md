---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_users"
description: |-
  Use this data source to get the list of DMS RabbitMQ users.
---

# huaweicloud_dms_rabbitmq_users

Use this data source to get the list of DMS RabbitMQ users.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_rabbitmq_users" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the RabbitMQ instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - Indicates the users list.

  The [users](#users_struct) structure is documented below.

<a name="users_struct"></a>
The `users` block supports:

* `access_key` - Indicates the user name.

* `vhosts` - Indicates the virtual hosts to be granted permissions for.

  The [vhosts](#users_vhosts_struct) structure is documented below.

<a name="users_vhosts_struct"></a>
The `vhosts` block supports:

* `write` - Indicates the granting resource write permissions using regular expressions.

* `read` - Indicates the granting resource read permissions using regular expressions.

* `vhost` - Indicates the name of the virtual host to be granted permissions for.

* `conf` - Indicates the granting resource permissions using regular expressions.
