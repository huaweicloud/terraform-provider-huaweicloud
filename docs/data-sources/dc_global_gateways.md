---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_global_gateways"
description: |-
  Use this data source to get the list of DC global gateways.
---

# huaweicloud_dc_global_gateways

Use this data source to get the list of DC global gateways.

## Example Usage

```hcl
data "huaweicloud_dc_global_gateways" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fields` - (Optional, List) Specifies the list of fields to be displayed. If omitted, all fields will be queried.

* `sort_key` - (Optional, String) Specifies the sorting field. Defaults to **id**.

* `sort_dir` - (Optional, String) Specifies the sorting order of returned results.
  There are two options: **asc** (default) and **desc**.

* `global_gateway_ids` - (Optional, List) Specifies the resource IDs for querying instances.

* `names` - (Optional, List) Specifies the resource names for querying instances.

* `enterprise_project_ids` - (Optional, List) Specifies the enterprise project IDs for querying instances.

* `site_network_ids` - (Optional, List) Specifies the site network IDs.

* `cloud_connection_ids` - (Optional, List) Specifies the cloud connection IDs.

* `statuses` - (Optional, List) Specifies the statuses by which instances are filtered.

* `global_center_network_ids` - (Optional, List) Specifies the central network IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `gateways` - The list of the global DC gateways.

  The [gateways](#gateways_struct) structure is documented below.

<a name="gateways_struct"></a>
The `gateways` block supports:

* `locales` - The locale address description information.

  The [locales](#gateways_locales_struct) structure is documented below.

* `current_peer_link_count` - The number of peer links allowed on a global DC gateway, indicating the number of
  enterprise routers that the global DC gateway can be attached to.

* `tags` - The key/value pairs to associate with the DC global gateway.

* `name` - The name of the global DC gateway.

* `location_name` - The location where the underlying device of the global DC gateway is deployed.

* `status` - The status of the global DC gateway.

* `description` - The description of the global DC gateway.

* `enterprise_project_id` - The enterprise project ID that the global DC gateway belongs to.

* `global_center_network_id` - The ID of the central network that the global DC gateway is added to.

* `bgp_asn` - The BGP ASN of the global DC gateway.

* `address_family` - The IP address family of the global DC gateway.
  + **ipv4**: Only IPv4 is supported.
  + **dual**: Both IPv4 and IPv6 are supported.

* `id` - The global DC gateway ID.

* `available_peer_link_count` - The number of peer links that can be created for a global DC gateway.

* `created_time` - The time when the global DC gateway was created.

* `updated_time` - The time when the global DC gateway was updated.

* `reason` - The cause of the failure to create the global DC gateway.

<a name="gateways_locales_struct"></a>
The `locales` block supports:

* `en_us` - The region name in English.

* `zh_cn` - The region name in Chinese.
