---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_unblock_records"
description: |-
  Use this data source to get the list of AAD unblock records within HuaweiCloud.
---

# huaweicloud_aad_unblock_records

Use this data source to get the list of AAD unblock records within HuaweiCloud.

## Example Usage

```hcl
variable "domain_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_aad_unblock_records" "test" {
  domain_id  = var.domain_id
  start_time = var.start_time
  end_time   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, String) Specified the account ID of IAM user.

* `start_time` - (Required, Int) Specified the start time of unblock record.

* `end_time` - (Required, Int) Specified the end time of unblock record.

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `unblock_record` - The unblock record.
  The [unblock_record](#unblock_record_struct) structure is documented below.

<a name="unblock_record_struct"></a>
The `unblock_record` block supports:

* `ip` - The IP address.

* `executor` - The executor.

* `block_id` - The block id.

* `blocking_time` - The blocking time, the value is a timestamp.

* `unblocking_time` - The unblocking time, the value is a timestamp.

* `unblock_type` - The unblock type. The valid values are as follows:
  + **manual**: Indicates manual unblock.
  + **automatic**: Indicates automatic unblock.

* `status` - The unblock status. The valid values are as follows:
  + **unblocking**: Indicates unblocking status.
  + **success**: Indicates successful status.
  + **failed**: Indicates failed status.

* `sort_time` - The sort time, the value is a timestamp.
