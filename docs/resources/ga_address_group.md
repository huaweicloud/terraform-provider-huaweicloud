---
subcategory: "Global Accelerator (GA)"
---

# huaweicloud_ga_address_group

Manages a GA IP address group resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

resource "huaweicloud_ga_address_group" "test" {
  name        = var.name
  description = "Created by terraform"

  ip_addresses {
    cidr        = "192.168.1.0/24"
    description = "The IP addresses included in the address group"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the IP address group name.

* `description` - (Optional, String) Specifies the description of the IP address group.

* `ip_addresses` - (Optional, List, ForceNew) Specifies the list of CIDR block configurations of the IP address group.
  The [ip_addresses](#address_group_ip_addresses) structure is documented below.

<a name="address_group_ip_addresses"></a>
The `ip_addresses` block supports:

* `cidr` - (Required, String, ForceNew) Specifies the CIDR block associated with the IP address group.

* `description` - (Optional, String, ForceNew) Specifies the description of the associated CIDR block.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the IP address group.
  + **ACTIVE**: The resource is running.

* `created_at` - The creation time of the IP address group.

* `updated_at` - The lasted update time of the IP address group.

* `ip_addresses` - The list of CIDR block configurations of the IP address group.
  The [ip_addresses](#address_group_ip_addresses_attr) structure is documented below.

* `listeners` - The listener associated with the IP address group.
  The [listeners](#address_group_associated_listeners) structure is documented below.

<a name="address_group_ip_addresses_attr"></a>
The `ip_addresses` block supports:

* `created_at` - The creation time of the CIDR block associated with the IP address group.

<a name="address_group_associated_listeners"></a>
The `listeners` block supports:

* `id` - The ID of the listener associated with the IP address group.

* `type` - The listener type associated with the IP address group.
  + **BLACK**: The blacklsit.
  + **WHITE**: The whitelist.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.
* `delete` - Default is 3 minutes.

## Import

The IP address group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ga_address_group.test <id>
```
