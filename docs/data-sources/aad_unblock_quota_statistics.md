---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_unblock_quota_statistics"
description: |-
  Use this data source to get information about the AAD unblock quota statistics within HuaweiCloud.
---

# huaweicloud_aad_unblock_quota_statistics

Use this data source to get information about the AAD unblock quota statistics within HuaweiCloud.

## Example Usage

```hcl
variable "domain_id" {}

data "huaweicloud_aad_unblock_quota_statistics" "test" {
  domain_id = var.domain_id
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, String) Specified the account ID of IAM user.

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `type` - The user type. The valid values are as follows:
  + **common_user**: Indicates common user.
  + **native_protection_user**: Indicates native basic protection user.

* `total_unblocking_quota` - The total unblocking quota.

* `remaining_unblocking_quota` - The remaining unblocking quota.

* `unblocking_quota_today` - The unblocking quota of today.

* `remaining_unblocking_quota_today` - The remaining unblocking quota of today.
