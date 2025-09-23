---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_block_statistics"
description: |-
  Use this data source to get information about the AAD block statistics within HuaweiCloud.
---

# huaweicloud_aad_block_statistics

Use this data source to get information about the AAD block statistics within HuaweiCloud.

## Example Usage

```hcl
variable "domain_id" {}

data "huaweicloud_aad_block_statistics" "test" {
  domain_id = var.domain_id
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, String) Specified the account ID of IAM user.

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `total_unblocking_times` - The total unblocking times.

* `manual_unblocking_times` - The manual unblocking times.

* `automatic_unblocking_times` - The automatic unblocking times.

* `current_blocked_ip_numbers` - The current blocked IP number.
