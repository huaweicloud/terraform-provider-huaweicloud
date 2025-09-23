---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_listener"
description: |-
  Manages a GA listener resource within HuaweiCloud.
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

* `frozen_info` - The frozen details of cloud services or resources.
  The [frozen_info](#Listener_frozen_info) structure is documented below.

<a name="Listener_frozen_info"></a>
The `frozen_info` block supports:

* `status` - The status of a cloud service or resource.
  The valid values are as follows:
  + `0`: unfrozen/normal (The cloud service will recover after being unfrozen.)
  + `1`: frozen (Resources and data will be retained, but the cloud service cannot be used.)
  + `2`: deleted/terminated (Both resources and data will be cleared.)

* `effect` - The status of the resource after being forzen.
  The valid values are as follows:
  + `1` (default): The resource is frozen and can be released.
  + `2`: The resource is frozen and cannot be released.
  + `3`: The resource is frozen and cannot be renewed.

* `scene` - The service scenario.
  The valid values are as follows:
  + **ARREAR**: The cloud service is in arrears, including expiration of yearly/monthly resources and fee deduction
    failure of pay-per-use resources.
  + **POLICE**: The cloud service is frozen for public security.
  + **ILLEGAL**: The cloud service is frozen due to violation of laws and regulations.
  + **VERIFY**: The cloud service is frozen because the user fails to pass the real-name authentication.
  + **PARTNER**: A partner freezes their customer's resources.

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
