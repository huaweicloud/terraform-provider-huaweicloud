---
subcategory: "Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_antiddos_quota"
description: |-
  Use this data source to query the quota information of Anti-DDoS service.
---

# huaweicloud_antiddos_quota

Use this data source to query the quota information of Anti-DDoS service.

## Example Usage

```hcl
data "huaweicloud_antiddos_quota" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `current` - The current used quota.

* `quota` - The maximum quota limit.
