---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_connection_numbers"
description: |-
  Use this data source to get the list of Advanced Anti-DDos connection numbers within HuaweiCloud.
---

# huaweicloud_aad_connection_numbers

Use this data source to get the list of Advanced Anti-DDos connection numbers within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}
variable "ip" {}

data "huaweicloud_aad_connection_numbers" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
  end_time    = var.end_time
  ip          = var.ip
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies the instance ID.

* `start_time` - (Required, String) Specifies the start time.

* `end_time` - (Required, String) Specifies the end time.

* `ip` - (Required, String) Specifies the IP address.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The connection number data list.  
  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `name` - The connection number name.

* `list` - The connection number data items.  
  The [list](#list_struct) structure is documented below.

<a name="list_struct"></a>
The `list` block supports:

* `time` - The timestamp in milliseconds.

* `value` - The connection number value.
