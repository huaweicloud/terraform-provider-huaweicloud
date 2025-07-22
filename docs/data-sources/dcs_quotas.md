---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_quotas"
description: |-
  Use this data source to get the default instance quota and total memory quota of a tenant and the maximum and minimum
  quotas a tenant can apply for.
---

# huaweicloud_dcs_quotas

Use this data source to get the default instance quota and total memory quota of a tenant and the maximum and minimum
quotas a tenant can apply for.

## Example Usage

```hcl
data "huaweicloud_dcs_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the quota information.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - Indicates the list of quotas.

  The [resources](#quotas_resources_struct) structure is documented below.

<a name="quotas_resources_struct"></a>
The `resources` block supports:

* `used` - Indicates the number of created instances and used memory.

* `type` - Indicates the quota type.
  The value can be **instance** or **ram**.
  + **instances**: instance quota
  + **ram**: memory quota

* `unit` - Indicates the resource unit.
  + When `type` is **instance**, no value is returned.
  + When `type` is **ram**, **GB** is returned.

* `min` - Indicates the minimum limit.
  + It is minimum limit of instance quota when `type` is **instance**.
  + It is minimum limit of memory quota when `type` is **ram**.

* `max` - Indicates the maximum limit.
  + It is maximum limit of instance quota when `type` is **instance**.
  + It is maximum limit of memory quota when `type` is **ram**.

* `quota` - Indicates the max number of instances that can be created and max allowed total memory.
