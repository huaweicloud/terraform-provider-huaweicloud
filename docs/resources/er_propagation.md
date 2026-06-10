---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_propagation"
description: |-
  Manages a propagation resource under the route table for ER service within HuaweiCloud.
---

# huaweicloud_er_propagation

Manages a propagation resource under the route table for ER service within HuaweiCloud.

## Example Usage

### Create a propagation to a route table using a VPC attachment

```hcl
variable "instance_id" {}
variable "route_table_id" {}
variable "vpc_attachment_id" {}

resource "huaweicloud_er_propagation" "test" {
  instance_id    = var.instance_id
  route_table_id = var.route_table_id
  attachment_id  = var.vpc_attachment_id
}
```

### Create a propagation to a route table using a VPN Gateway attachment and with a custom route policy

```hcl
variable "instance_id" {}
variable "route_table_id" {}
variable "vpn_gateway_attachment_id" {}
variable "import_policy_id" {}

resource "huaweicloud_er_propagation" "with_route_policy" {
  instance_id    = var.instance_id
  route_table_id = var.route_table_id
  attachment_id  = var.vpn_gateway_attachment_id

  route_policy {
    import_policy_id = var.import_policy_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the ER instance and route table are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the ER instance to which the route table and the
  attachment belongs.

* `route_table_id` - (Required, String, NonUpdatable) Specifies the ID of the route table to which the propagation
  belongs.

* `attachment_id` - (Required, String, NonUpdatable) Specifies the ID of the attachment corresponding to the
  propagation.

* `route_policy` - (Optional, List) Specifies the import route policy configuration.  
  The [route_policy](#er_propagation_route_policy) structure is documented below.

  -> This parameter currently only applies to certain types of attachments, such as VPN gateway.<br>
     For information regarding support for more attachment types, please contact the relevant service via a ticket.

<a name="er_propagation_route_policy"></a>
The `route_policy` block supports:

* `import_policy_id` - (Optional, String) Specifies the import route policy ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `attachment_type` - The type of the attachment corresponding to the propagation.
  + **vpc**: Virtual Private Cloud
  + **vpn**: VPN Gateway
  + **vgw**: Virtual Gateway for Cloud Dedicated Line
  + **peering**: Peer-to-peer connection (creating peer-to-peer connections by loading enterprise routers in different
    regions via Cloud Connect CC)
  + **gdgw**: Global Access Gateway
  + **cfw**: Cloud Firewall

* `status` - The current status of the propagation.

* `created_at` - The creation time.

* `updated_at` - The latest update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 2 minutes.

## Import

Propagations can be imported using their `id` and the related `instance_id` and `route_table_id`, separated by
slashes (/), e.g.

```bash
$ terraform import huaweicloud_er_propagation.test <instance_id>/<route_table_id>/<id>
```
