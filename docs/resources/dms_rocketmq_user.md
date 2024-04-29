---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_user"
description: ""
---

# huaweicloud_dms_rocketmq_user

Manages DMS RocketMQ user resources within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_rocketmq_user" "test" {
  instance_id          = var.instance_id
  access_key           = "user_test"
  secret_key           = "abcdefg"
  white_remote_address = "10.10.10.10"
  admin                = false
  default_topic_perm   = "PUB"
  default_group_perm   = "PUB"
  
  topic_perms {
    name = "topic_name"
    perm = "PUB"
  }
  
  group_perms {
    name = "group_name"
    perm = "PUB"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the rocketMQ instance.
  Changing this parameter will create a new resource.

* `access_key` - (Required, String, ForceNew) Specifies the name of the user, which starts with a letter, consists of 7
  to 64 characters and can contain only letters, digits, hyphens (-), and underscores (_).
  Changing this parameter will create a new resource.

* `secret_key` - (Required, String, ForceNew) Specifies the password of the user. Use 8 to 32 characters. Contain at
  least three of the following character types:
  + Uppercase letters.
  + Lowercase letters.
  + Digits.
  + Special characters \`~!@#$%^&*()-_=+\|[{}];:'"",<.>/?. Cannot be the `access_key` or the `access_key` spelled
    backwards.
  
  Changing this parameter will create a new resource.

* `white_remote_address` - (Optional, String) Specifies the IP address whitelist.

* `admin` - (Optional, Bool) Specifies whether the user is an administrator.

* `default_topic_perm` - (Optional, String) Specifies the default topic permissions.
  Value options: **PUB|SUB**, **PUB**, **SUB**, **DENY**.

* `default_group_perm` - (Optional, String) Specifies the default consumer group permissions.
  Value options: **PUB|SUB**, **PUB**, **SUB**, **DENY**.

* `topic_perms` - (Optional, List) Specifies the special topic permissions.
  The [permission](#DmsRocketMQUser_PermsRef) structure is documented below.

* `group_perms` - (Optional, List) Specifies the special consumer group permissions.
  The [permission](#DmsRocketMQUser_PermsRef) structure is documented below.

<a name="DmsRocketMQUser_PermsRef"></a>
The `topic_perms` and `group_perms` block supports:

* `name` - (Optional, String) Indicates the name of a topic or consumer group.

* `perm` - (Optional, String) Indicates the permissions of the topic or consumer group.
  Value options: **PUB|SUB**, **PUB**, **SUB**, **DENY**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The rocketmq user can be imported using the rocketMQ `instance_id` and user `access_key` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dms_rocketmq_user.test <instance_id>/<access_key>
```
