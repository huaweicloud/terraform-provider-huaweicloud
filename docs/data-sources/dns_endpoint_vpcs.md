---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_endpoint_vpcs"
description: |-
  Use this data source to get the list of DNS endpoint VPCs.
---

# huaweicloud_dns_endpoint_vpcs

Use this data source to get the list of DNS endpoint VPCs.

## Example Usage

```hcl
data "huaweicloud_dns_endpoint_vpcs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `vpcs` - Indicates the VPCs of an endpoint.

  The [vpcs](#vpcs_struct) structure is documented below.

<a name="vpcs_struct"></a>
The `vpcs` block supports:

* `id` - Indicates the VPC ID, which is a UUID used to identify the VPC.

* `inbound_endpoint_count` - Indicates the number of inbound endpoints in a VPC.

* `outbound_endpoint_count` - Indicates the number of outbound endpoints in a VPC.
