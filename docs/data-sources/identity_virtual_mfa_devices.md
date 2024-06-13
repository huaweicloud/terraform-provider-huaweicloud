---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_virtual_mfa_devices"
description: |-
  Use this data source to get the list of IAM virtual MFA devices within HuaweiCloud.
---

# huaweicloud_identity_virtual_mfa_devices

Use this data source to get the list of IAM virtual MFA devices within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identity_virtual_mfa_devices" "test" {}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Optional, String) Specifies the user ID to which the virtual MFA device belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `virtual_mfa_devices` - The list of virtual MFA devices.

  The [virtual_mfa_devices](#virtual_mfa_devices_struct) structure is documented below.

<a name="virtual_mfa_devices_struct"></a>
The `virtual_mfa_devices` block supports:

* `serial_number` - The virtual MFA device serial number.

* `user_id` - The user ID to which the virtual MFA device belongs.
