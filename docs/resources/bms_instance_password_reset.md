---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_instance_password_reset"
description: |-
  Manages a BMS instance password reset resource within HuaweiCloud.
---

# huaweicloud_bms_instance_password_reset

Manages a BMS instance password reset resource within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}
variable "new_password" {}

resource "huaweicloud_bms_instance_password_reset" "test" {
  server_id    = var.server_id
  new_password = var.new_password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the provider-level
  region will be used. Changing this creates a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the BMS ID.

* `new_password` - (Required, String, NonUpdatable) Specifies the new BMS password. The password must comply with the
  following rules:
  + Consists of 8 to 26 characters.
  + Contains at least three of the following character types:
      - Uppercase letters
      - Lowercase letters
      - Digits
      - Special characters:
          * Windows: !@$%-_=+[]:./?
          * Linux: !@%^-_=+[]{}:,./?
  + Cannot contain the username or the username in reverse.
  + Cannot contain more than two characters in the same sequence as they appear in the username. (This requirement
    applies only to Windows BMSs.)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the BMS instance ID.
