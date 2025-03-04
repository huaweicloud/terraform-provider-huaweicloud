---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_quotas"
description: |-
  Use this data source to get a list of VPC resource quotas.
---

# huaweicloud_vpc_quotas

Use this data source to get a list of VPC resource quotas.

## Example Usage

```hcl
data "huaweicloud_vpc_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the type of resource to filter quotas.
  The value can be **vpc**, **subnet**, **securityGroup**, **securityGroupRule**, **publicIp**,
  **vpn**, **vpngw**, **vpcPeer**, **firewall**, **shareBandwidth**, **shareBandwidthIP**,
  **loadbalancer**, **listener**, **physicalConnect**, **virtualInterface**,
  **vpcContainRoutetable**, and **routetableContainRoutes**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The quota objects.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - The resource objects.

  The [resources](#quotas_resources_struct) structure is documented below.

<a name="quotas_resources_struct"></a>
The `resources` block supports:

* `min` - The minimum quota value allowed.

* `type` - The type of resource to filter quotas.

* `used` - The number of created resources.

* `quota` - The maximum quota values for the resources.
