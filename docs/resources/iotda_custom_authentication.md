---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_custom_authentication"
description: |-
  Manages an IoTDA custom authentication resource within HuaweiCloud.
---

# huaweicloud_iotda_custom_authentication

Manages an IoTDA custom authentication resource within HuaweiCloud.

-> When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify the IoTDA service
  endpoint in `provider` block.
  You can login to the IoTDA console, choose the instance **Overview** and click **Access Details**
  to view the HTTPS application access address. An example of the access address might be
  **9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com**, then you need to configure the
  `provider` block as follows:

  ```hcl
  provider "huaweicloud" {
    endpoints = {
      iotda = "https://9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com"
    }
  }
  ```

## Example Usage

```hcl
variable "name" {}
variable "func_urn" {}
variable "signing_token" {}
variable "signing_public_key" {}

resource "huaweicloud_iotda_custom_authentication" "test" {
  authorizer_name    = var.name
  func_urn           = var.func_urn
  signing_token      = var.signing_token
  signing_public_key = var.signing_public_key
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the custom authentication resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `authorizer_name` - (Required, String) Specifies the name of the custom authentication.
  The name contains a maximum of `128` characters, and only letters, digits, underscores (_), and hyphens (-)
  are allowed. The name must be unique.

* `func_urn` - (Required, String) Specifies the URN of the function associated with the custom authentication.

* `signing_enable` - (Optional, Bool) Specifies whether to enable signature authentication. Defaults to **true**.
  You are advised to enable this function. If this function is enabled, authentication information that does not
  meet signature requirements will be rejected to reduce invalid function calls.

* `signing_token` - (Optional, String) Specifies the private key for signature authentication.
  The key contains a maximum of `128` characters, and only letters, digits, underscores (_), and hyphens (-)
  are allowed.

* `signing_public_key` - (Optional, String) Specifies the public key for signature authentication.
  Used to check whether the signature information carried by the device is correct.

  -> 1. The parameters `signing_token` and `signing_public_key` are mandatory when `signing_enable` is set to **true**.
  <br/>2. The parameter `signing_public_key` must be RSA encryption public key.

* `default_authorizer` - (Optional, Bool) Specifies whether the custom authentication is the default
  authentication mode. Defaults to **false**.
  If this parameter is set to **true**, the current authentication policy is used fo authentication on all devices that
  support SNI unless otherwise specified.

* `status` - (Optional, String) Specifies whether to enable the custom authentication mode. Defaults to **INACTIVE**.
  The valid values are as follows:
  + **ACTIVE**: The authentication is enabled.
  + **INACTIVE**: The authentication is disabled.

* `cache_enable` - (Optional, Bool) Specifies whether to enable the cache function. Defaults to **false**.
  If this parameter is set to **true** and the device input parameters (username, client ID, password, certificate
  information, and function URN) remain unchanged, the cache result is directly used when the cache result exists.
  Yor are advised to set this parameter to **false** during debugging, set this parameter to **true** during production
  to avoid frequent function invoking.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `func_name` - The name of the function associated with the custom authentication.

* `create_time` - The creation time of the custom authentication.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.

* `update_time` - The latest update time of the custom authentication.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.

## Import

The custom authentication can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_custom_authentication.test <id>
```
