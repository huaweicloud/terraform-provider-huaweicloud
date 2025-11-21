---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user_token_info"
description: |-
  Use this data source to get the list of virtual MFA devices in the Identity and Access Management V5 service.
---

# huaweicloud_identityv5_virtual_mfa_devices

Use this data source to get the list of virtual MFA devices in the Identity and Access Management V5 service.

## Example Usage

```hcl
data "huaweicloud_identityV5_virtual_mfa_devices" "devices" {}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Optional, String) Specifies the ID of the IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `devices` - Indicates the list of virtual MFA devices. Structure is documented below.
  The [devices](#IdentityMFA_Devices) structure is documented below.

<a name="IdentityMFA_Devices"></a>
The `devices` block supports:

* `enabled` - Indicates whether the MFA device is enabled.

* `serial_number` - Indicates the serial number of the MFA device.

* `user_id` - Indicates the ID of the IAM user.
