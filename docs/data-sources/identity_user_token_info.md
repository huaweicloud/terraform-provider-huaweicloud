---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user_token_info"
description: "Use this data-source to get information of an IAM user token within HuaweiCloud."
---

# huaweicloud_identity_user_token_info

Use this data-source to get information of an IAM user token within HuaweiCloud.

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

data "huaweicloud_identity_user_token_info" "test" {
  token = huaweicloud_identity_user_token.test.token
  no_catalog = "false"
}
```

## Argument Reference

* `token` - (Required, String) Specifies the token to be verified.

* `nocatalog` - (Optional, String) Specifies whether to include the catalog field in the response. The value can be:
  + **true**: not include;
  + **false**: include

  Defaults to **false**.

## Attribute Reference

* `catalog` - service catalog information, including the services and endpoints the user has access to. This field will
  not be returned if **nocatalog** parameter is set to **true**.
  The [catalog](#IdentityToken_Catalog) structure is documented below.

* `domain` - Indicates the domain information of the IAM user if the scope is set to domain.
  The [domain](#IdentityToken_Domain) structure is documented below.

* `expires_at` - Indicates the time when the token will expire. The value is a UTC time in the
  **YYYY-MM-DDTHH:mm:ss.ssssssZ** format.

* `issued_at` - Indicates the time when the token was issued. The value is a UTC time in the
  **YYYY-MM-DDTHH:mm:ss.ssssssZ** format.

* `methods` - Indicates the methods used to obtain the token, such as **password**, **token**, etc.

* `project` - Indicates the project information of the IAM user if the scope is set to **project**.
  The [project](#IdentityToken_Project) structure is documented below.

* `roles` - Indicates the roles and permissions associated with the token.
  The [roles](#IdentityToken_Roles) structure is documented below.

* `user` -Indicates the detailed information of the IAM user associated with the token.
  The [user](#IdentityToken_User) structure is documented below.

<a name="IdentityToken_Catalog"></a>
The `catalog` block contains:

* `endpoint` - Indicates service catalog information, including the services and endpoints the user has access to. This
  field will not be returned if **nocatalog** parameter is set.
  The [endpoint](#IdentityToken_Endpoint) structure is documented below.

* `id` - Indicates the ID of service

* `name` - Indicates the name of service

* `type` - Indicates the Interface type of service

<a name="IdentityToken_Endpoint"></a>
The `endpoint` block contains:

* `id` - Indicates the ID of Endpoint.

* `interface` - Indicates Interface type, describing the visibility of the interface at this endpoint. A value of
  'public' indicates that this interface is a public interface.

* `region` - Indicates the region of endpoint.

* `region_id` - Indicates the region ID of endpoint.

* `url` - Indicates the url of endpoint.

<a name="IdentityToken_Domain"></a>
The `domain` block contains:

* `id` - Indicates the ID of domain

* `name` - Indicates the name of domain

<a name="IdentityToken_Project"></a>
The `project` block contains:

* `domain` - Indicates the information of domain.
  The [domain](#IdentityToken_Project_Domain) structure is documented below.

* `id` - Indicates the ID of domain.

* `name` - Indicates the name of domain.

<a name="IdentityToken_Project_Domain"></a>
The `domain` block contains:

* `id` - Indicates the ID of domain.

* `name` - Indicates the name of domain.

<a name="IdentityToken_Roles"></a>
The `roles` block contains:

* `id` - Indicates the ID of domain.

* `name` - Indicates the name of domain.

<a name="IdentityToken_User"></a>
The `user` block contains:

* `domain` - Indicates the information of domain.
  The [domain](#IdentityToken_User_Domain) structure is documented below.

* `id` - Indicates the ID of domain.

* `name` - Indicates the name of domain.

* `password_expires_at` - Indicates The password expiration time, "" indicates that the password does not expire.
  Note: UTC time, formatted as YYYY-MM-DDTHH:mm:ss.ssssssZ, date and timestamp format refers to ISO-8601,
  for example: 2023-06-28T08:56:33.710000Z.

<a name="IdentityToken_User_Domain"></a>
The `domain` block contains:

* `id` - Indicates the ID of IAM domain.

* `name` - Indicates the name of IAM domain.
