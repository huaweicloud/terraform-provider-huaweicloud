---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_block_list"
description: |-
  Use this data source to get the list of blocked IPs in CNAD Advanced.
---

# huaweicloud_cnad_advanced_block_list

Use this data source to get the list of blocked IPs in CNAD Advanced.

## Example Usage

```hcl
data "huaweicloud_cnad_advanced_block_list" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `blocking_list` - The list of blocked IPs.

  The [blocking_list](#advancedBlockList) structure is documented below.

<a name="advancedBlockList"></a>
The `blocking_list` block supports:

* `ip` - The blocked IP address.

* `blocking_time` - The timestamp when the IP was blocked.

* `estimated_unblocking_time` - The estimated timestamp when the IP will be unblocked.

* `status` - The status of the blocked IP. Possible values are:
  + **unblocking**: The IP is being unblocked.
  + **success**: The IP has been successfully unblocked.
  + **failed**: The unblocking operation failed.
