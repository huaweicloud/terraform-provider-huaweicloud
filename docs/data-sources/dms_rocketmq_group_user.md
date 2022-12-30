---
subcategory: "Distributed Message Service (DMS)"
---

# huaweicloud_dms_group_user

Use this data source to get the list of users that have been granted permissions for a consumer group.

## Example Usage

```hcl
variable "instance_id" {}
variable "group" {}

data "huaweicloud_dms_rocketmq_group_user" "test" {
  instance_id = var.instance_id
  group       = var.group
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

* `group` - (Required, String) Specifies the name of the RocketMQ consumer group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `policies` - Indicates the list of user associated with the consumer group.
  The [Policy](#DmsRocketMQGroupUser_Policy) structure is documented below.

<a name="DmsRocketMQGroupUser_Policy"></a>
The `Policy` block supports:

* `access_key` - Indicates the access key of the user.

* `secret_key` - Indicates the secret key of the user.

* `white_remote_address` - Indicates the IP address whitelist.

* `admin` - Indicates whether the user is an administrator.

* `perm` - Indicates the permissions of the user.
