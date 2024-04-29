---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_address_groups"
description: ""
---

# huaweicloud_ga_address_groups

Use this data source to get the list of IP address groups.

## Example Usage

```hcl
variable "associated_listener_id" {}

data "huaweicloud_ga_address_groups" "test" {
  listener_id = var.associated_listener_id
}
```

## Argument Reference

The following arguments are supported:

* `address_group_id` - (Optional, String) Specifies the ID of the IP address group.

* `name` - (Optional, String) Specifies the name of the IP address group.

* `status` - (Optional, String) Specifies the status of the IP address group.
  The valid values are as follows:
  + **ACTIVE**: The status of the IP address group is normal operation.
  + **ERROR**: The status of the IP address group is error.

* `listener_id` - (Optional, String) Specifies the ID of the listener associated with the IP address group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `address_groups` - The list of the IP address groups.
  The [address_groups](#address_groups_address_group) structure is documented below.

<a name="address_groups_address_group"></a>
The `address_groups` block supports:

* `id` - The ID of the IP address group.

* `name` - The name of the IP address group.  

* `description` - The description of the IP address group.

* `status` - The status of the IP address group.

* `ip_addresses` - The list of CIDR block configurations of the IP address group.
  The [ip_addresses](#address_groups_ip_list) structure is documented below.

* `associated_listeners` - The list of the listeners associated with the IP address group.
  The [associated_listeners](#address_groups_associated_listeners) structure is documented below.

* `created_at` - The creation time of the IP address group.

* `updated_at` - The latest update time of the IP address group.

<a name="address_groups_ip_list"></a>
The `ip_addresses` block supports:

* `cidr` - The CIDR block included in the IP address group.

* `description` - The description of the CIDR block.

* `created_at` - The creation time of the CIDR block.

<a name="address_groups_associated_listeners"></a>
The `associated_listeners` block supports:

* `id` - The ID of the listener associated with the IP address group.

* `type` - The listener type associated with the IP address group.
  The value can be one of the following:
  + **BLACK**: The blacklsit.
  + **WHITE**: The whitelist.
