---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_quotas"
description: |-
  Use this data source to get the AS quota list within HuaweiCloud.
---

# huaweicloud_as_quotas

Use this data source to get the AS quota list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_as_quotas" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The quota details.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - The quota resources.

  The [resources](#quotas_resources_struct) structure is documented below.

<a name="quotas_resources_struct"></a>
The `resources` block supports:

* `type` - The quota type. Values are:
  + **scaling_Group**: AS group quota.
  + **scaling_Config**: AS configuration quota.
  + **scaling_Policy**: AS policy quota.
  + **scaling_Instance**: Instance quota.
  + **bandwidth_scaling_policy**: Bandwidth scaling policy quota.

* `min` - The quota lower limit.

* `used` - The used amount of the quota.

* `quota` - The total quota.

* `max` - The quota upper limit.
