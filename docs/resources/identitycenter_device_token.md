---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_device_authorization"
description: ""
---

# huaweicloud_identitycenter_device_authorization

Manages an Identity Center device authorization resource within HuaweiCloud.

## Example Usage

```hcl
variable "client_id" {}
variable "client_secret" {}
variable "device_code" {}
resource "huaweicloud_identitycenter_device_authorization" "test"{
	client_id       = var.client_id
	client_secret   = var.client_secret
    device_code		= var.device_code
    grant_type		= "urn:ietf:params:oauth:grant-type:device_code"
}
```

## Argument Reference

The following arguments are supported:

* `client_id` - (Required, String, ForceNew) Unique ID of the client registered in the IAM Identity Center.

  Changing this parameter will create a new resource.

* `client_secret` - (Required, String, ForceNew) Secret string generated for the client to obtain authorization from services in subsequent calls.

  Changing this parameter will create a new resource.

* `code` - (Optional, String, ForceNew) Authorization code received from the authorization service. This parameter is required when executing an authorization request to obtain access to the token.

  Changing this parameter will create a new resource.

* `device_code` - (Optional, String, ForceNew) Used only when the authorization type (grant_type) is the device code (urn:ietf:params:oauth:grant-type:device_code).

  Changing this parameter will create a new resource.

* `grant_type` - (Required, String, ForceNew) Authorization type, which can be authorization code, device code.

  Changing this parameter will create a new resource.

* `redirect_uri` - (Optional, String, ForceNew) Application URL that will receive the authorization code. The user authorizes a service to send a request to this URL.

  Changing this parameter will create a new resource.

* `refresh_token` - (Optional, String, ForceNew) Refresh token, which can be used to obtain a new access token after the original access token expires.

  Changing this parameter will create a new resource.

* `scopes` - (Optional, List, ForceNew) List of scopes defined by a client to restrict permissions for access token authorization.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `access_token` - Opaque token used to access IAM Identity Center resources assigned to users 

* `expires_in` - Expiration time (in seconds) of an access token.

* `id_token` - Opaque token used to identify a user.

* `refresh_token` - Refresh token, which can be used to obtain a new access token after the original access token expires.