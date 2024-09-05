---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_group_quotas"
description: |-
  Use this data source to get the quota list of a specified AS group within HuaweiCloud.
---

# huaweicloud_as_group_quotas

Use this data source to get the quota list of a specified AS group within HuaweiCloud.

## Example Usage

```hcl
variable "scaling_group_id" {}

data "huaweicloud_as_group_quotas" "test" {
  scaling_group_id = var.scaling_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `scaling_group_id` - (Required, String) Specifies the AS group ID to query quotas.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The quota details.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - The quota resource list.

  The [resources](#quotas_resources_struct) structure is documented below.

<a name="quotas_resources_struct"></a>
The `resources` block supports:

* `type` - The quota type. Valid values are:
  + **scaling_Policy**: Indicates AS policies.
  + **scaling_Instance**: Indicates instances.

* `max` - The quota upper limit.

* `min` - The quota lower limit.

* `used` - The used quota.

* `quota` - The total quota.
