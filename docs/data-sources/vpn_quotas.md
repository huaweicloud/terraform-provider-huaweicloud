---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_quotas"
description: |-
  Use this data source to get the list of VPN quotas.
---

# huaweicloud_vpn_quotas

Use this data source to get the list of VPN quotas.

## Example Usage

```hcl
data "huaweicloud_vpn_quotas" "test" {}
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

* `type` - The resource type.

* `quota` - The quota upper limit.

* `used` - The number of resources in use.
