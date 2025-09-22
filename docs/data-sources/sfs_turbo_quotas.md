---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_quotas"
description: |-
  Use this data source to get the list of SFS Turbo quotas and their usage.
---

# huaweicloud_sfs_turbo_quotas

Use this data source to get the list of SFS Turbo quotas and their usage.

## Example Usage

```hcl
data "huaweicloud_sfs_turbo_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of SFS Turbo quotas and their usage.
  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block contains:

* `resources` - The list of quota resources.
  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block contains:

* `type` - The type of the quota. The value can be **shares** or **capacity**.
* `max` - The maximum value of the quota.
* `min` - The minimum value of the quota.
* `quota` - The total quota value.
* `unit` - The unit of the quota.
* `used` - The used quota value.
