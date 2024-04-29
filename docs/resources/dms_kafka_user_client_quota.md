---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_user_client_quota"
description: ""
---

# huaweicloud_dms_kafka_user_client_quota

Manage DMS kafka user client quota resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_kafka_user_client_quota" "test" {
  instance_id        = var.instance_id
  user               = "test_user"
  client             = "consumer_group"
  producer_byte_rate = 2048
  consumer_byte_rate = 1024
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the kafka instance.
  Changing this parameter will create a new resource.

* `user` - (Optional, String, ForceNew) Specifies the user name to apply the quota. It must be empty
  if the value of `user_default` is **true**. Changing this parameter will create a new resource.

* `user_default` - (Optional, Bool, ForceNew) Specifies the user default configuration of the quota.
  If `user_default` is **true**, the quota applies to all users. It can not be **true** if the value of `user`
  is not empty. Changing this parameter will create a new resource.

* `client` - (Optional, String, ForceNew) Specifies the ID of the client to which the quota applies.
  It must be empty if the value of `client_default` is **true**. Changing this parameter will create a new resource.

* `client_default` - (Optional, Bool, ForceNew) Specifies the client default configuration of the quota.
  If `client_default` is **true**, the quota applies to all clients. It can not be **true** if the value of
  `client` is not empty. Changing this parameter will create a new resource.

-> **NOTE:** At least one of `user`, `user_default`, `client` and `client_default` must be specified.

* `producer_byte_rate` - (Optional, Int) Specifies an upper limit on the prodction rate. The unit is B/s.
  If this parameter is left blank, no limit is set.

* `consumer_byte_rate` - (Optional, Int) Specifies an upper limit on the consumption rate. The unit is B/s.
  If this parameter is left blank, no limit is set.

-> **NOTE:** At least one of `producer_byte_rate` and `consumer_byte_rate` must be specified.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 50 minutes.
* `update` - Default is 50 minutes.
* `delete` - Default is 50 minutes.

## Import

The kafka user client quota can be imported using the kafka `instance_id`, `user`, `user_default`, `client` and `client_default`
separated by slashes, e.g.

```bash
$ terraform import huaweicloud_dms_kafka_user_client_quota.test <instance_id>/<user>/<user_default>/<client>/<client_default>
```
