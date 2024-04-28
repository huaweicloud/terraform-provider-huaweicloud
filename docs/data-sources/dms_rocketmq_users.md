---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_users"
description: ""
---

# huaweicloud_dms_rocketmq_users

Use this data source to get the list of DMS rocketMQ users.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_dms_rocketmq_users" "test" {
  instance_id          = var.instance_id
  access_key           = "user0001"
  white_remote_address = "10.10.10.10"
  admin                = false
  default_topic_perm   = "PUB"
  default_group_perm   = "SUB"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the rocketMQ instance.

* `access_key` - (Optional, String) Specifies the user name.

* `white_remote_address` - (Optional, String) Specifies the IP address whitelist.

* `admin` - (Optional, Bool) Specifies whether the user is an administrator.

* `default_topic_perm` - (Optional, String) Specifies the default topic permissions.
  Value options: **PUB|SUB**, **PUB**, **SUB**, **DENY**.

* `default_group_perm` - (Optional, String) Specifies the default consumer group permissions.
  Value options: **SUB**, **DENY**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - The list of users.
  The [users](#DMS_rockermq_users) structure is documented below.

<a name="DMS_rockermq_users"></a>
The `users` block supports:

* `access_key` - Indicates the name of the user.

* `white_remote_address` - Indicates the IP address whitelist.

* `admin` - Indicates whether the user is an administrator.

* `default_topic_perm` - Indicates the default topic permissions.
  Value options: **PUB|SUB**, **PUB**, **SUB**, **DENY**.

* `default_group_perm` - Indicates the default consumer group permissions.
  Value options: **SUB**, **DENY**.

* `topic_perms` - The list of the special topic permissions.
  The [topic_perms](#DMS_rocketmq_users_topic_perms) structure is documented below.

* `group_perms` - The list of the special consumer group permissions.
  The [group_perms](#DMS_rocketmq_users_group_perms) structure is documented below.

<a name="DMS_rocketmq_users_topic_perms"></a>
The `topic_perms` block supports:

* `name` - Indicates the name of a topic.

* `perm` - Indicates the permissions of the topic.
  Value options: **PUB|SUB**, **PUB**, **SUB**, **DENY**.

<a name="DMS_rocketmq_users_group_perms"></a>
The `group_perms` block supports:

* `name` - Indicates the name of consumer group.

* `perm` - Indicates the permissions of consumer group.
  Value options: **SUB**, **DENY**.
