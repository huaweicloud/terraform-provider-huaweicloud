---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_users"
description: |-
  Use this data source to get the list of Kafka instance users.
---

# huaweicloud_dms_kafka_users

Use this data source to get the list of Kafka instance users.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_users" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `name` - (Optional, String) Specifies the user name.

* `description` - (Optional, String) Specifies the user description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - Indicates the user list.

  The [users](#users_struct) structure is documented below.

<a name="users_struct"></a>
The `users` block supports:

* `name` - Indicates the username.

* `description` - Indicates the description.

* `role` - Indicates the user role.

* `default_app` - Indicates whether the application is the default application.

* `created_at` - Indicates the create time.
