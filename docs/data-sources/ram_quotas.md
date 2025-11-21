---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_quotas"
description: |-
  Use this data source to list quotas in Resource Access Manager.
---

# huaweicloud_ram_quotas

Use this data source to list quotas in Resource Access Manager.

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

  The [quotas](#quotas) structure is documented below.

<a name="quotas"></a>
The `quotas` block supports:

* `resources` - The list of resources.

  The [resources](#resources) structure is documented below.

<a name="resources"></a>
The `resources` block supports:

* `type` - The type of resource quota.

* `quota` - The available quota of the resource.

* `min` - The minimum quota of the resource.

* `max` - The maximum quota of the resource.

* `used` - The used quota of the resource.
