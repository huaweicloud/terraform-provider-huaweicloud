---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_token"
description: ""
---

# huaweicloud_identity_token

Verify the validity of the token data source within Huaweicloud

-> **NOTE:** You *must* have admin privileges to use this data source.

## Example Usage

```hcl
variable "account_name" {}
variable "user_name" {}
variable "password" {}

resource "huaweicloud_identity_user_token" "test" {
  account_name = var.account_name
  user_name    = var.user_name
  password     = var.password
}

data "huaweicloud_identity_token" "test" {
  token = huaweicloud_identity_user_token.test.token
}
```

## Argument Reference

* `nocatalog` - (Optional, String) If set, the response will not include the **catalog** field. Any non-empty string will be treated as **true**.

## Attribute Reference

* `catalog` - Service catalog information, including the services and endpoints the user has access to. This field will not be returned if **nocatalog** parameter is set.

* `domain` - The domain information of the IAM user if the scope is set to domain.

* `expires_at` - The time when the token will expire. The value is a UTC time in the YYYY-MM-DDTHH:mm:ss.ssssssZ format.

* `issued_at` - The time when the token was issued. The value is a UTC time in the YYYY-MM-DDTHH:mm:ss.ssssssZ format.

* `methods` - The methods used to obtain the token, such as **password**, **token**, etc.

* `project` - The project information of the IAM user if the scope is set to **project**.

* `roles` - The roles and permissions associated with the token.

* `user` - The detailed information of the IAM user associated with the token.

