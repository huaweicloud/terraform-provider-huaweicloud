---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_global_gateway_peer_links"
description: |-
  Use this data source to get the list of DC global gateway peer links.
---

# huaweicloud_dc_global_gateway_peer_links

Use this data source to get the list of DC global gateway peer links.

## Example Usage

```hcl
variable "global_dc_gateway_id" {}

data "huaweicloud_dc_global_gateway_peer_links" "test" {
  global_dc_gateway_id = var.global_dc_gateway_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `global_dc_gateway_id` - (Required, String) Specifies the global DC gateway ID.

* `fields` - (Optional, List) Specifies the list of fields to be displayed. If omitted, all fields will be queried.

* `sort_key` - (Optional, String) Specifies the sorting field. Defaults to **id**.

* `sort_dir` - (Optional, String) Specifies the sorting order of returned results.
  There are two options: **asc** (default) and **desc**.

* `peer_link_ids` - (Optional, List) Specifies the resource IDs for querying instances.

* `names` - (Optional, List) Specifies the resource names for querying instances.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `peer_links` - The list of the peer links.

  The [peer_links](#peer_links_struct) structure is documented below.

<a name="peer_links_struct"></a>
The `peer_links` block supports:

* `id` - The peer link ID.

* `name` - The name of the peer link.

* `description` - The description of the peer link.

* `reason` - The cause of the failure to add the peer link.

* `global_dc_gateway_id` - The ID of the global DC gateway that the peer link is added for.

* `bandwidth_info` - The bandwidth information.
  The [bandwidth_info](#peer_link_bandwidth_info) structure is documented below.

* `peer_site` - The bandwidth information.
  The [peer_site](#peer_link_peer_site) structure is documented below.

* `status` - The status of the peer link. This attribute values include:
  + **PENDING_CREATE**: The peer link is being created.
  + **PENDING_UPDATE**: The peer link is being updated.
  + **ACTIVE**: The peer link is available.
  + **ERROR**: An error occurred.

* `created_time` - The time when the peer link was added.

* `updated_time` - The time when the peer link was updated.

* `create_owner` - The cloud service where the peer link is used. This attribute values include:
  + **cc**: Cloud Connect.
  + **dc**: Direct Connect.

* `instance_id` - The ID of the instance associated with the peer link.

<a name="peer_link_bandwidth_info"></a>
The `bandwidth_info` block supports:

* `bandwidth_size` - The bandwidth size.

* `gcb_id` - The global connection bandwidth ID.

<a name="peer_link_peer_site"></a>
The `peer_site` block supports:

* `gateway_id` - The ID of enterprise router that the global DC gateway is attached to.

* `project_id` - The project ID of the enterprise router that the global DC gateway is attached to.

* `region_id` - The region ID of the enterprise router that the global DC gateway is attached to.

* `link_id` - The connection ID of the peer gateway at the peer site.
  For example, if the peer gateway is an enterprise router, this attribute means attachment ID.
  If the peer gateway is a global DC gateway, this attribute means the peer link ID.

* `site_code` - The site information of the global DC gateway.

* `type` - The type of the peer gateway. This attribute values include:
  + **ER**: Enterprise router.
  + **GDGW**: Global DC gateway.
