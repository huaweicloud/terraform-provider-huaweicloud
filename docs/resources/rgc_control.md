---
subcategory: "RGC"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_control"
description: |-
  Manages an RGC control resource within HuaweiCloud.
---

# huaweicloud_rgc_control

Manages an RGC control resource within HuaweiCloud.

## Example Usage

```hcl
variable "identifier" {}
variable "target_identifier" {}
variable "parameters" {}

resource "huaweicloud_rgc_control" "test" {
  identifier        = var.identifier
  target_identifier = var.target_identifier
  parameters        = var.parameters
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `identifier` - (Required, String, NonUpdatable) Specifies the id of the enabled control.

* `target_identifier` - (Required, String, NonUpdatable) Specifies the id of the organizational unit.

* `parameters` - (Optional, String, NonUpdatable) Specifies the parameter information for control strategies with parameters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `state` - Indicates the state of enabled control.

* `version` - Indicates the version of enabled control.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 15 minutes.

* `delete` - Default is 15 minutes.

## Import

The RGC control can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rgc_control.test <id>
```
