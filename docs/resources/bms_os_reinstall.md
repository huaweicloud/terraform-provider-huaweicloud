---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_os_reinstall"
description: |-
  Manages a BMS OS reinstall resource within HuaweiCloud.
---

# huaweicloud_bms_os_reinstall

Manages a BMS OS reinstall resource within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}
variable "key_name" {}
variable "user_id" {}

resource "huaweicloud_bms_os_reinstall" "test" {
  server_id = var.server_id

  os_reinstall{
    key_name = var.key_name
    user_id  = var.user_id

    metadata {
      __system__encrypted = "0"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the BMS ID.

* `os_reinstall` - (Required, List, NonUpdatable) Specifies the operation of reinstalling the BMS OS.
  The [os_reinstall](#os_reinstall_struct) structure is documented below.

<a name="os_reinstall_struct"></a>
The `os_reinstall` block supports:

* `admin_pass` - (Optional, String, NonUpdatable) Specifies the initial password of the BMS administrator account. The
  Linux administrator is **root**, and the Windows administrator is **Administrator**. Recommended password complexity
  requirements are as follows:
  + The password contains 8 to 26 characters.
  + Contains at least three of the following character types: uppercase letters, lowercase letters, digits, and special
    characters **!@$%^-_=+[{}]:,./?**.
  + The password cannot contain the username or the username in reverse.

-> **NOTE:**
  1.For Windows BMSs, the password cannot contain more than two consecutive characters in the username.
  <br/>2.For Linux BMSs, `user_data` can be used to inject a password. In this case, `admin_pass` is invalid.
  <br/>3.Either `admin_pass` or `key_name` can be set.
  <br/>4.If both `admin_pass` and `key_name` are empty, `user_data` in metadata must be set.

* `key_name` - (Optional, String, NonUpdatable) Specifies the key pair name.

* `user_id` - (Optional, String, NonUpdatable) Specifies the user ID.

* `metadata` - (Optional, List, NonUpdatable) Specifies the BMS metadata.
  The [metadata](#os_reinstall_metadata_struct) structure is documented below.

<a name="os_reinstall_metadata_struct"></a>
The `metadata` block supports:

* `user_data` - (Optional, String, NonUpdatable) Specifies the Linux image root password injected during the BMS OS
  reinstallation. It is a user-defined initial password. The password change script must be encoded using Base64.
  Recommended password complexity requirements are as follows:
  + Contains 8 to 26 characters.
  + Contains at least three of the following character types: uppercase letters, lowercase letters, digits, and special
    **characters !@$%^-_=+[{}]:,./?**.

* `__system__encrypted` - (Optional, String, NonUpdatable) Specifies whether the system disk is encrypted. The value can
  be **0 (not encrypted)** or **1 (encrypted)**.
  + If this parameter is not specified, the system disk will not be encrypted by default.
  + If this parameter is set to **0**, `__system__cmkid` will be invalid.

* `__system__cmkid` - (Optional, String, NonUpdatable) Specifies the CMK ID which is used for encryption. This parameter
  is used with `__system__encrypted`.

* `__system__encryption_algorithm` - (Optional, String, NonUpdatable) Specifies the encryption algorithms for encrypted
  volumes. Value options: **AES_256**, **SM4**. Defaults to **AES_256**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
