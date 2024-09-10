---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_endpoint_group"
description: ""
---

# huaweicloud_ga_endpoint_group

Manages a GA endpoint group resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "description" {}
variable "listener_id" {}

resource "huaweicloud_ga_endpoint_group" "test" {
  name        = var.name
  description = var.description
  region_id   = "cn-south-1"

  listeners {
    id = var.listener_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the endpoint group name.  
  The name can contain `1` to `64` characters, only letters, digits, and hyphens (-) are allowed.

* `region_id` - (Required, String, ForceNew) Specifies the region where the endpoint group belongs.

  Changing this parameter will create a new resource.

* `listeners` - (Required, List, ForceNew) Specifies the listeners associated with the endpoint group.
  The [Id](#EndpointGroup_Id) structure is documented below.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the information about the endpoint group.  
  The description contain a maximum of `255` characters, and the angle brackets (< and >) are not allowed.

* `traffic_dial_percentage` - (Optional, Int) Specifies the percentage of traffic distributed to the endpoint group.
  The value ranges from `0` to `100`. Defaults to `100`.

<a name="EndpointGroup_Id"></a>
The `Id` block supports:

* `id` - (Required, String) Specifies the ID of the associated listener.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the provisioning status. The value can be one of the following:
  + **ACTIVE**: The resource is running.
  + **PENDING**: The status is to be determined.
  + **ERROR**: Failed to create the resource.
  + **DELETING**: The resource is being deleted.

* `created_at` - Indicates when the endpoint group was created.

* `updated_at` - Indicates when the endpoint group was updated.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The endpoint group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ga_endpoint_group.test <id>
```
