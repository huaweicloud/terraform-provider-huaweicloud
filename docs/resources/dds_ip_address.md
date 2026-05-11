---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_ip_address"
description: |-
  Manages a shard/config IP address resource within HuaweiCloud.
---

# huaweicloud_dds_lts_log

Manages a shard/config IP address resource within HuaweiCloud.

-> 1. The frozen cluster DDS instance does not support add shard/config IP address.
  <br/>2. The cluster DDS instance associated with the IPv6 subnet do not support add shard/config IP address.

## Example Usage

```hcl
variable "instance_id" {}
variable "type" {}
variable "password" {}

resource "huaweicloud_dds_ip_address" "test" {
  instance_id = var.instance_id
  type        = var.type
  password    = var.password
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster DDS instance.

* `type` - (Required, String, NonUpdatable) Specifies the node type of the cluster DDS instance.
  The value can be **shard** or **config**.

* `password` - (Required, String, NonUpdatable) Specifies the password for enabling this function
  for a cluster DDS instance.
  The value must be `8` to `32` characters in length and contain uppercase letters, lowercase letters,
  digits and special characters, such as **~!@#%^*-_=+?**.

* `target_ids` - (Optional, List) Specifies the shard group IDs.
  This parameter only supports in update, creating resource does not support.
  This parameter only supports add the shard group ID, does not support remove the shard group ID.
  If the shard or config IP address is added for the first time, leave this parameter empty.
  If a shard IP address has been added to a DDS instance, you need to specify this parameter to add an IP address
  to the new shard group.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID, also is the DDS instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 10 minutes.
