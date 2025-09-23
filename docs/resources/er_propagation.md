---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_propagation"
description: ""
---

# huaweicloud_er_propagation

Manages a propagation resource under the route table for ER service within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "route_table_id" {}
variable "attachment_id" {}

resource "huaweicloud_er_propagation" "test" {
  instance_id    = var.instance_id
  route_table_id = var.route_table_id
  attachment_id  = var.attachment_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the ER instance and route table are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the ER instance to which the route table and the
  attachment belongs.  
  Changing this parameter will create a new resource.

* `route_table_id` - (Required, String, ForceNew) Specifies the ID of the route table to which the propagation
  belongs.  
  Changing this parameter will create a new resource.

* `attachment_id` - (Required, String, ForceNew) Specifies the ID of the attachment corresponding to the propagation.  
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `attachment_type` - The type of the attachment corresponding to the propagation.

* `status` - The current status of the propagation.

* `created_at` - The creation time.

* `updated_at` - The latest update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 2 minutes.

## Import

Propagations can be imported using their `id` and the related `instance_id` and `route_table_id`, separated by
slashes (/), e.g.

```bash
$ terraform import huaweicloud_er_propagation.test <instance_id>/<route_table_id>/<id>
```
