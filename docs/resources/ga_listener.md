---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_listener"
description: ""
---

# huaweicloud_ga_listener

Manages a GA listener resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "description" {}
variable "accelerator_id" {}

resource "huaweicloud_ga_listener" "test" {
  accelerator_id = var.accelerator_id
  name           = var.name
  description    = var.description
  protocol       = "TCP"

  port_ranges {
    from_port = 4000
    to_port   = 4200
  }
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, String, ForceNew) Specifies the ID of the global accelerator associated with the listener.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the listener name.  
  The name can contain `1` to `64` characters, only letters, digits, and hyphens (-) are allowed.

* `port_ranges` - (Required, List) Specifies the port range used by the listener.
  The [PortRange](#Listener_PortRange) structure is documented below.

* `protocol` - (Required, String, ForceNew) Specifies the protocol used by the listener to forward requests.
  The value can be **TCP** or **UDP**.

  Changing this parameter will create a new resource.

* `client_affinity` - (Optional, String) Specifies the client affinity. The value can be one of the following:
  + **Source IP address**: Requests from the same IP address are forwarded to the same endpoint.
  + **NONE**: Requests are evenly forwarded across the endpoints.

  Defaults to **NONE**.

* `description` - (Optional, String) Specifies the information about the listener.  
  The description contain a maximum of `255` characters, and the angle brackets (< and >) are not allowed.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the listener.

<a name="Listener_PortRange"></a>
The `PortRange` block supports:

* `from_port` - (Required, Int) Specifies the start port number.
  The valid value is range from `1` to `65,535`.

* `to_port` - (Required, Int) Specifies the end port number.
  The valid value is range from `1` to `65,535`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the provisioning status. The value can be one of the following:
  + **ACTIVE**: The resource is running.
  + **PENDING**: The status is to be determined.
  + **ERROR**: Failed to create the resource.
  + **DELETING**: The resource is being deleted.

* `created_at` - Indicates when the listener was created.

* `updated_at` - Indicates when the listener was updated.
  
## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The listener can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ga_listener.test <id>
```
