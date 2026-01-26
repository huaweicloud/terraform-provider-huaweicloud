---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_quotas"
description: |-
  Use this data source to get the list of quotas in Resource Access Manager.
---

# huaweicloud_ram_quotas

Use this data source to get the list of quotas in Resource Access Manager.

## Example Usage

```hcl
data "huaweicloud_ram_quotas" "test" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of quotas.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - The list of resources.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `type` - The type of quota.

* `quota` - The total number quotas.

* `min` - The minimum quota.

* `max` - The maximum quota.

* `used` - The number of quotas already used
