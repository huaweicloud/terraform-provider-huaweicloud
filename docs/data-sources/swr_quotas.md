---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_quotas"
description: |-
  Use this data source to get the list of SWR quotas.
---

# huaweicloud_swr_quotas

Use this data source to get the list of SWR quotas.

## Example Usage

```hcl
data "huaweicloud_swr_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of quotas.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `quota_limit` - The quota limit.

* `quota_key` - The quota type.

* `unit` - The quota unit.

* `used` - The used quota.
