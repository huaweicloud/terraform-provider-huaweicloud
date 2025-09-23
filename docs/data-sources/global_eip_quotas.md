---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_quotas"
description: |-
  Use this data source to get the list of GEIP resource quotas.
---

# huaweicloud_global_eip_quotas

Use this data source to get the list of GEIP resource quotas.

## Example Usage

```hcl
data "huaweicloud_global_eip_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, List) Specifies the resource type.
  Valid values are **geip**, **geip_segment**, **internetBandwidthIP**, **internetBandwidth**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - Indicates the resources list.

  The [resources](#quotas_resources_struct) structure is documented below.

<a name="quotas_resources_struct"></a>
The `resources` block supports:

* `type` - Indicates the quota type.

* `used` - Indicates the used num.

* `quota` - Indicates the total quotas.

* `min` - Indicates the min quotas.
