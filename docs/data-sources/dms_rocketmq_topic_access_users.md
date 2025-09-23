---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_topic_access_users"
description: |-
  Use this data source to get the list of RocketMQ topic access users.
---

# huaweicloud_dms_rocketmq_topic_access_users

Use this data source to get the list of RocketMQ topic access users.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic" {}

data "huaweicloud_dms_rocketmq_topic_access_users" "test" {
  instance_id = var.instance_id
  topic       = var.topic
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `topic` - (Required, String) Specifies the topic name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - Indicates the user list.

  The [policies](#policies_struct) structure is documented below.

<a name="policies_struct"></a>
The `policies` block supports:

* `admin` - Indicates whether the user is an administrator.

* `perm` - Indicates the permissions.

* `access_key` - Indicates the user name.

* `white_remote_address` - Indicates the IP address whitelist.
