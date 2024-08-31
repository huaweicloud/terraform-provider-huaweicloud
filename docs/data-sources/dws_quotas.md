---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_quotas"
description: |-
  Use this data source to query the list of available resource quotas within HuaweiCloud.
---

# huaweicloud_dws_quotas

Use this data source to query the list of available resource quotas within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dws_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - All quotas that match the filter parameters.

  The [quotas](#quotas_quotas_struct) structure is documented below.

<a name="quotas_quotas_struct"></a>
The `quotas` block supports:

* `type` - The type of the quota.

* `used` - The number of quotas used.

* `limit` - The number of available quotas.

* `unit` - The unit of the quota.
