---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_consumer_group_access_users"
description: |-
  Use this data source to get the list of RocketMQ consumer group access users.
---

# huaweicloud_dms_rocketmq_consumer_group_access_users

Use this data source to get the list of RocketMQ consumer group access users.

## Example Usage

```hcl
variable "instance_id" {}
variable "group" {}

data "huaweicloud_dms_rocketmq_consumer_group_access_users" "test" {
  instance_id = var.instance_id
  group       = var.group
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `group` - (Required, String) Specifies the consumer group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - Indicates the user list.

  The [policies](#policies_struct) structure is documented below.

<a name="policies_struct"></a>
The `policies` block supports:

* `white_remote_address` - Indicates the IP address whitelist.

* `admin` - Indicates whether the user is an administrator.

* `perm` - Indicates the permissions.

* `access_key` - Indicates the user name.
