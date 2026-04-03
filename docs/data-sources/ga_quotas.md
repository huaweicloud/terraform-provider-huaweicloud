---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_quotas"
description: |-
Use this data source to get the list of GA quotas within HuaweiCloud.
---

# huaweicloud_ga_quotas

Use this data source to get the list of GA quotas within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_ga_quotas" "test" {}
```

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `quotas` - All quotas that match the filter parameters.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `type` - The quota mark.

* `min` - The minimum quota threshold.

* `max` - The maximum quota threshold.

* `quota` - The quota size.
