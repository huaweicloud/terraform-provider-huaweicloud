---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_hosted_connects"
description: |-
  Use this data source to get a list of hosted connects.
---

# huaweicloud_dc_hosted_connects

Use this data source to get a list of hosted connects.

## Example Usage

```hcl
data "huaweicloud_dc_hosted_connects" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `hosting_id` - (Optional, List) Specifies the hosting direct connect IDs to which the hosted connects belong.

* `hosted_connect_id` - (Optional, List) Specifies the hosted connect IDs.

* `name` - (Optional, List) Specifies the hosted connect names.

* `sort_key` - (Optional, String) Specifies the sorting field.

* `sort_dir` - (Optional, List) Specifies the sorting order of returned results.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `hosted_connects` - Indicates the list of hosted connects.
  The [hosted_connects](#HostedConnects_HostedConnect) structure is documented below.

<a name="HostedConnects_HostedConnect"></a>
The `hosted_connects` block supports:

* `id` - Indicates the hosted connect ID.

* `name` - Indicates the hosted connect name.

* `description` - Indicates the hosted connect description.

* `hosting_id` - Indicates the hosting direct connect ID.

* `vlan` - Indicates the VLAN allocated to the hosted connect.

* `bandwidth` - Indicates the bandwidth of the hosted connect in Mbit/s.

* `provider` - Indicates the provider of the hosted connect.

* `provider_status` - Indicates the provider status of the hosted connect.

* `type` - Indicates the type of the hosted connect.

* `port_type` - Indicates the port type of the hosted connect.

* `location` - Indicates the location of the hosted connect.

* `peer_location` - Indicates the peer location of the hosted connect.

* `status` - Indicates the status of the hosted connect.

* `apply_time` - Indicates the time when the hosted connect was applied.

* `create_time` - Indicates the time when the hosted connect was created.
