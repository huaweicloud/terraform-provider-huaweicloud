---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_quotas"
description: |-
  Use this data source to get the list of IMS quotas within HuaweiCloud.
---

# huaweicloud_ims_quotas

Use this data source to get the list of IMS quotas within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_ims_quotas" "test" {}
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

* `type` - The resource type. The valid value is **image**.

* `used` - The number of resource quotas already in use.

* `quota` - The total quota of resources.

* `min` - The minimum quota of resources.

* `max` - The maximum quota of resources.
