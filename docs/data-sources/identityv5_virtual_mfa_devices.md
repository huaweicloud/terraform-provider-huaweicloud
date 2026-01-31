---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityV5_virtual_mfa_devices"
description: |-
  Use this data source to get the list of virtual MFA devices within HuaweiCloud.
---

# huaweicloud_identityv5_virtual_mfa_devices

Use this data source to get the list of virtual MFA devices within HuaweiCloud.

## Example Usage

### Query all virtual MFA devices

```hcl
data "huaweicloud_identityV5_virtual_mfa_devices" "test" {}
```

### Query virtual MFA devices by user ID

```hcl
variable "user_id" {}

data "huaweicloud_identityV5_virtual_mfa_devices" "test" {
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Optional, String) Specifies the ID of the IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `devices` - The list of virtual MFA devices that matched filter parameters.  
  The [devices](#v5_virtual_mfa_devices) structure is documented below.

<a name="v5_virtual_mfa_devices"></a>
The `devices` block supports:

* `enabled` - Whether the MFA device is enabled.

* `serial_number` - The serial number of the MFA device.

* `user_id` - The ID of the IAM user.
