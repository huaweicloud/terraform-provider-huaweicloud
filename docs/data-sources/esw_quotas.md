---
subcategory: "Enterprise Switch (ESW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_esw_quotas"
description: |-
  Use this data source to get the ESW instance quotas.
---

# huaweicloud_esw_quotas

Use this data source to get the ESW instance quotas.

## Example Usage

```hcl
data "huaweicloud_esw_quotas" "test" {}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the quotas. If omitted, the provider-level region will be
  used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the ESW instance quotas.
  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - Indicates the list of quota infos.
  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `type` - Indicates the type of the quota.

* `quota` - Indicates the total quota.

* `used` - Indicates the used quota.

* `min` - Indicates the minimum total quota.

* `max` - Indicates the maximum total quota.
