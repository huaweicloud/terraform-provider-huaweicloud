---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_instance_quota"
description: |-
  Use this data source to get CBH instance quota within HuaweiCloud.
---

# huaweicloud_cbh_instance_quota

Use this data source to get CBH instance quota within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cbh_instance_quota" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the CBH instance quota.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `quota` - The maximum number of CBH instances that can be created.

* `quota_used` - The current number of CBH instances created.
