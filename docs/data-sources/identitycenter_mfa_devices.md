---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_mfa_devices"
description: |-
  Use this data source to get the Identity Center user mfa devices.
---

# huaweicloud_identitycenter_mfa_devices

Use this data source to get the Identity Center user mfa devices.

## Example Usage

```hcl
variable "identity_store_id" {}
variable "user_id" {}

data "huaweicloud_identitycenter_mfa_devices" "test"{
  identity_store_id = var.identity_store_id
  user_id           = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `identity_store_id` - (Required, String) Specifies the ID of the identity store that associated with IAM Identity
  Center.

* `user_id` - (Required, String) Specifies the ID of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `mfa_devices` - The mfa devices of the service provider.
  The [mfa_devices](#mfa_devices_struct) structure is documented below.

<a name="mfa_devices_struct"></a>
The `mfa_devices` block supports:

* `device_id` - The ID of the mfa device.

* `device_name` - The name of the mfa device.

* `display_name` - The display name of the mfa device.

* `mfa_type` - The type of the mfa device.

* `registered_date` - The registered date of the mfa device.
