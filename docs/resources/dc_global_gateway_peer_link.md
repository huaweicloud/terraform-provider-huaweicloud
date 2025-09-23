---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_global_gateway_peer_link"
description: |-
  Manages a DC global gateway peer link resource within HuaweiCloud.
---

# huaweicloud_dc_global_gateway_peer_link

Manages a DC global gateway peer link resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "global_dc_gateway_id" {}
variable "gateway_id" {}
variable "project_id" {}
variable "region_id" {}

resource "huaweicloud_dc_global_gateway_peer_link" "test" {
  name                 = var.name
  global_dc_gateway_id = var.global_dc_gateway_id
  description          = "test description"

  peer_site {
    gateway_id = var.gateway_id
    project_id = var.project_id
    region_id  = var.region_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `name` - (Required, String) Specifies the name of the global DC gateway peer link.

* `global_dc_gateway_id` - (Required, String, NonUpdatable) Specifies the global DC gateway ID.

  -> It is required that the gateway has created a virtual interface.

* `peer_site` - (Required, List, NonUpdatable) Specifies the site of the peer link.
  Currently, only one site information can be configured.
  The [peer_site](#peer_link_peer_site) structure is documented below.

* `description` - (Optional, String) Specifies the description of the global DC gateway peer link.

<a name="peer_link_peer_site"></a>
The `peer_site` block supports:

* `gateway_id` - (Required, String, NonUpdatable) Specifies the ID of enterprise router (ER) that the global DC gateway
  is attached to.

* `project_id` - (Required, String, NonUpdatable) Specifies the project ID of the enterprise router (ER) that the global
  DC gateway is attached to.

* `region_id` - (Required, String, NonUpdatable) Specifies the region ID of the enterprise router (ER) that the global
  DC gateway is attached to.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `reason` - The cause of the failure to add the peer link.

* `bandwidth_info` - The bandwidth information.
  The [bandwidth_info](#peer_link_bandwidth_info) structure is documented below.

* `peer_site/link_id` - The connection ID of the peer gateway at the peer site.
  For example, if the peer gateway is an enterprise router, this attribute means attachment ID.
  If the peer gateway is a global DC gateway, this attribute means the peer link ID.

* `peer_site/site_code` - The site information of the global DC gateway.

* `peer_site/type` - The type of the peer gateway. This attribute values include:
  + **ER**: Enterprise router.
  + **GDGW**: Global DC gateway.

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

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Defaults to 10 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 10 minutes.

## Import

The DC global gateway peer link resource can be imported using the `global_dc_gateway_id` and `id`,
separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dc_global_gateway_peer_link.test <global_dc_gateway_id>/<id>
```
