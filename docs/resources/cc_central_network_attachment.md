---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_network_attachment"
description: ""
---

# huaweicloud_cc_central_network_attachment

Manages a central network attachment resource of Cloud Connect within HuaweiCloud.

To allow network instances such as enterprise routers and global DC gateways to communicate with each other
across regions, you need to add these network instances to the central network.

## Example Usage

```hcl
variable "central_network_id" {}
variable "enterprise_router_id" {}
variable "enterprise_router_project_id" {}
variable "enterprise_router_region_id" {}
variable "global_dc_gateway_id" {}
variable "global_dc_gateway_project_id" {}
variable "global_dc_gateway_region_id" {}

resource "huaweicloud_cc_central_network_attachment" "test" {
  name                         = "demo"
  description                  = "This is a demo"
  central_network_id           = var.central_network_id
  enterprise_router_id         = var.enterprise_router_id
  enterprise_router_project_id = var.enterprise_router_project_id
  enterprise_router_region_id  = var.enterprise_router_region_id
  global_dc_gateway_id         = var.global_dc_gateway_id
  global_dc_gateway_project_id = var.global_dc_gateway_project_id
  global_dc_gateway_region_id  = var.global_dc_gateway_region_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `central_network_id` - (Required, String, ForceNew) The central network ID.

  Changing this parameter will create a new resource.

* `name` - (Required, String) The name of the attachment.

* `enterprise_router_id` - (Required, String, ForceNew) The enterprise router ID.

  Changing this parameter will create a new resource.

* `enterprise_router_project_id` - (Required, String, ForceNew) The project ID to which the enterprise router belongs.

  Changing this parameter will create a new resource.

* `enterprise_router_region_id` - (Required, String, ForceNew) The region ID to which the enterprise router belongs.

  Changing this parameter will create a new resource.

* `global_dc_gateway_id` - (Required, String, ForceNew) The global DC gateway ID.

  Changing this parameter will create a new resource.

* `global_dc_gateway_project_id` - (Required, String, ForceNew) The project ID to which the global DC gateway belongs.

  Changing this parameter will create a new resource.

* `global_dc_gateway_region_id` - (Required, String, ForceNew) The region ID to which the global DC gateway belongs.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) The description of the attachment.

* `central_network_plane_id` - (Optional, String, ForceNew) The central network plane ID.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `state` - Central network attachment status.
  The valid values are as follows:
    - **AVAILABLE**
    - **CREATING**
    - **UPDATING**
    - **DELETING**
    - **FREEZING**
    - **UNFREEZING**
    - **RECOVERING**
    - **FAILED**
    - **DELETED**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The central network attachment can be imported using `central_network_id`, `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cc_central_network_attachment.test <central_network_id>/<id>
```
