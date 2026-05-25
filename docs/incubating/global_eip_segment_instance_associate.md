---
subcategory: "EIP (Elastic IP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_segment_instance_associate"
description: |-
  Manages a Global EIP segment instance associate resource within HuaweiCloud.
---

# huaweicloud_global_eip_segment_instance_associate

Manages a Global EIP segment instance associate resource within HuaweiCloud.

## Example Usage

```hcl
variable "global_eip_segment_id" {}
variable "region" {}
variable "instance_id" {}
variable "instance_type" {}
variable "project_id" {}

resource "huaweicloud_global_eip_segment_instance_associate" "test" {
  global_eip_segment_id = var.global_eip_segment_id
  
  global_eip_segment {
    region        = var.region
    instance_id   = var.instance_id
    instance_type = var.instance_type
    project_id    = var.project_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource is located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `global_eip_segment_id` - (Required, String, NonUpdatable) Specifies the ID of the global EIP segment.

* `global_eip_segment` - (Required, List, NonUpdatable) Specifies the configuration of the global EIP segment to
  associate with an instance. The [global_eip_segment](#block--global_eip_segment) structure is documented below.

<a name="block--global_eip_segment"></a>
The `global_eip_segment` block supports:

* `region` - (Required, String, NonUpdatable) Specifies the region where the instance is located.

* `instance_type` - (Required, String, NonUpdatable) Specifies the type of the instance to associate.
  Valid value is **DC-CONNECT-GATEWAY**.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the instance to associate.

* `project_id` - (Required, String, NonUpdatable) Specifies the project ID of the instance.

* `instance_site` - (Optional, String, NonUpdatable) Specifies the site ID of the instance.

* `service_id` - (Optional, String, NonUpdatable) Specifies the service ID of the instance.

* `service_type` - (Optional, String, NonUpdatable) Specifies the service type of the instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in the format of `<global_eip_segment_id>`.

* `associate_instance` - The information about the associated instance.
  The [associate_instance](#block--associate_instance) structure is documented below.

<a name="block--associate_instance"></a>
The `associate_instance` block exports:

* `instance_id` - The ID of the associated instance.

* `instance_type` - The type of the associated instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

This resource can be imported using the `global_eip_segment_id`, e.g.

```bash
$ terraform import huaweicloud_global_eip_segment_instance_associate.test <global_eip_segment_id>
```
