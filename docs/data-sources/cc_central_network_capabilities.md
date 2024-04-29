---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_network_capabilities"
description: ""
---

# huaweicloud_cc_central_network_capabilities

Use this data source to get the list of CC central network capabilities.

## Example Usage

```hcl
variable "capability" {}

data "huaweicloud_cc_central_network_capabilities" "test" {
  capability = var.capability
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `capability` - (Optional, String) Specifies the capability of the central network.
  The valid values are as follows:
  + **central-network.is-support**: Whether the central network is supported.
  + **central-network.is-support-enterprise-project**: Whether the central network supports enterprise projects.
  + **central-network.is-support-tag**: Whether the central network supports tags.
  + **connection-bandwidth.size-range**: The bandwidth range for a cross-site connection.
  + **connection-bandwidth.charge-mode**: The Billing mode of the global private bandwidth for assigning cross-site
  connection bandwidths.
  + **er-instance.support-regions**: The list of the regions where Enterprise Router is available.
  + **er-instance.support-ipv6-regions**: The list of the regions where Enterprise Router supports IPv6.
  + **er-instance.support-dscp-regions**: The list of the regions that support gold, silver, and bronze global private
  bandwidths.
  + **er-instance.support-sites**: The list of the sites where Enterprise Router is available.
  + **gdgw-attachment.is-support**: Whether global DC gateways as attachments are supported.
  + **gdgw-attachment.support-regions**: The list of the regions where global DC gateways are available.
  + **gdgw-attachment.support-sites**: The list of the sites where global DC gateways are available.
  + **er-route-table-attachment.is-support**: Whether The enterprise router route tables as attachments are supported.
  + **er-route-table-attachment.support-regions**: The list of regions where enterprise router route tables can be added
  as attachments.
  + **er-route-table-attachment.support-sites**: The list of sites where enterprise router route tables can be added as
  attachments.
  + **cloud-alliance.support-regions**: The list of regions that support Cloud Alliance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `capabilities` - Central network capability list.
  The [capabilities](#Capabilities) structure is documented below.

<a name="Capabilities"></a>
The `capabilities` block supports:

* `capability` - The capability of the central network.

* `domain_id` - The ID of the account that the central network belongs to.

* `specifications` - The specifications of the central network capability.
