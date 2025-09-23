---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_quotas"
description: |-
  Use this data source to get the list of CES quotas.
---

# huaweicloud_ces_quotas

Use this data source to get the list of CES quotas.

## Example Usage

```hcl
data "huaweicloud_ces_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The quota information.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - The resource quota list.

  The [resources](#quotas_resources_struct) structure is documented below.

<a name="quotas_resources_struct"></a>
The `resources` block supports:

* `type` - The quota type.

* `used` - The used amount of the quota.

* `quota` - The total amount of the quota.

* `unit` - The unit.
