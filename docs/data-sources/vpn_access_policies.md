---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_access_policies"
description: |-
  Use this data source to get the list of VPN access policies.
---

# huaweicloud_vpn_access_policies

Use this data source to get the list of VPN access policies.

## Example Usage

```hcl
variable "vpn_server_id" {}

data "huaweicloud_vpn_access_policies" "test" {
  vpn_server_id = var.vpn_server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `vpn_server_id` - (Required, String) Specifies the ID of a VPN server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `access_policies` - The VPN access policy list.

  The [access_policies](#access_policies_struct) structure is documented below.

<a name="access_policies_struct"></a>
The `access_policies` block supports:

* `id` - The ID of an access policy.

* `name` - The name of an access policy.

* `dest_ip_cidrs` - The destination CIDR block list.

* `user_group_id` - The ID of the associated user group.

* `user_group_name` - The name of the associated user group.

* `description` - The access policy description.

* `created_at` - The creation time.

* `updated_at` - The update time.
