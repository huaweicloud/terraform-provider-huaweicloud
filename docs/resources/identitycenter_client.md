---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_client"
description: ""
---

# huaweicloud_identitycenter_client

Manages an Identity Center client resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_identitycenter_client" "test"{
  client_name                = "client_test"
  client_type                = "public"
  token_endpoint_auth_method = "client_secret_post"
  grant_types 	             = ["urn:ietf:params:oauth:grant-type:device_code"]
  response_types             = ["code"]
  scopes		             = ["openid"]
}
```

## Argument Reference

The following arguments are supported:

* `client_name` - (Required, String, ForceNew) Client name.

  Changing this parameter will create a new resource.

* `client_type` - (Required, String, ForceNew) Client type. Only "public" is supported as a value.

  Changing this parameter will create a new resource.

* `token_endpoint_auth_method` - (Required, String, ForceNew) Authentication method required to send a request to the token endpoint. Only "client_secret_post" is supported as a value.

  Changing this parameter will create a new resource.

* `grant_types` - (Required, List, ForceNew) OAuth2.0 authorization type that a client can use at the token endpoint. Only "urn:ietf:params:oauth:grant-type:device_code" or "authorization_code" are supported as values.

  Changing this parameter will create a new resource.

* `response_types` - (Required, List, ForceNew) OAuth2.0 authorization type that a client can use at the authorization endpoint. Only "code" is supported as a value.

  Changing this parameter will create a new resource.

* `scopes` - (Optional, List, ForceNew) List of scopes defined by a client to restrict permissions for access token authorization.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `client_id` - Unique ID of a client application

* `client_secret` - Secret string generated for the client to obtain authorization from services in subsequent calls

* `client_secret_expires_at` - Expiration time of the client ID and secret key

