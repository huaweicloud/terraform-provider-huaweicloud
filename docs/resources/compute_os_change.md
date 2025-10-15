---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_os_change"
description: |-
  Manages an ECS OS change resource within HuaweiCloud.
---

# huaweicloud_compute_os_change

Manages an ECS OS change resource within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}

resource "huaweicloud_compute_os_change" "test" {
  cloud_init_installed = "false"
  server_id            = var.server_id

  os_change{
    imageid = data.huaweicloud_images_images.test.images[1].id
    userid  = "test"
    mode    = "withStopServer"

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

* `cloud_init_installed` - (Required, String, NonUpdatable) Specifies whether the image with Cloud-Init or Cloudbase-Init
  installed. Value options:
  + **true**: The image with Cloud-Init or Cloudbase-Init installed.
  + **false**: The image without Cloud-Init or Cloudbase-Init not installed.

* `server_id` - (Required, String, NonUpdatable) Specifies the ECS ID.

* `os_change` - (Required, List, NonUpdatable) Specifies the info of re-installs an ECS OS.
  The [os_change](#os_change_struct) structure is documented below.

<a name="os_change_struct"></a>
The `os_change` block supports:

* `imageid` - (Required, String, NonUpdatable) Specifies the ID of the new image in UUID format.

* `adminpass` - (Optional, String, NonUpdatable) Specifies the initial password of the ECS administrator. The Windows
  administrator username is **Administrator**, and the Linux administrator username is **root**. Constraints:
  + The Windows ECS password cannot contain the username, the username in reverse, or more than two characters in the same
    sequence as they appear in the username.
  + Linux ECSs can use `user_data` to inject passwords. In such a case, `adminpass` is unavailable.
  + Either `adminpass` or `keyname` is specified.
  + If both `adminpass` and `keyname` are empty, Linux ECSs can use `user_data` specified in `metadata`.
  + `adminpass`, `keyname`, and `user_data` in `metadata` can be empty only when a private image password is used or when
    a password is set after the OS is changed. The constraints are as follows:
      - Windows OSs do not support private image passwords.
      - If you intend to reset a password after the OS change, ensure that the **__os_feature_list** parameter of the image
        contains **{"onekey_resetpasswd": "true"}**.
  + The `password` must contain **8** to **26** characters.
  + The `password` must contain at least three of the following character types: uppercase letters, lowercase letters,
    digits, and special characters (!@$%^-_=+[{}]:,./?~#*).

* `keyname` - (Optional, String, NonUpdatable) Specifies the key name.

* `userid` - (Optional, String, NonUpdatable) Specifies the user ID. When the `keyname` parameter is being specified, the
  value of this parameter is used preferentially. If this parameter is left blank, the user ID in the token is used by default.

* `mode` - (Optional, String, NonUpdatable) Specifies whether the ECS supports OS changeation when the ECS is running.
  If the parameter value is **withStopServer**, the ECS supports OS changeation when the ECS is running. In such a case,
  the system automatically stops the ECS before changeing its OS.

* `metadata` - (Optional, List, NonUpdatable) Specifies the metadata of the ECS for which the OS is to be changed.
  The [metadata](#os_change_metadata_struct) structure is documented below.

<a name="os_change_metadata_struct"></a>
The `metadata` block supports:

* `__system__encrypted` - (Optional, String, NonUpdatable) Specifies the encryption field in `metadata`.
  + **0**: indicates a non-encrypted disk.
  + **1**: indicates an encrypted disk.

* `__system__cmkid` - (Optional, String, NonUpdatable) Specifies the CMK ID, which indicates encryption in `metadata`.
  This parameter must be used with `__system__encrypted`.

* `user_data` - (Optional, String, NonUpdatable) Specifies the user data to be injected to the ECS during the creation.
 Text and text files can be injected. It is only supported when `cloud_init_installed` is set to **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
