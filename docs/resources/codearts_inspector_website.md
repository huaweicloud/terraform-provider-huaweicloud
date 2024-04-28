---
subcategory: "CodeArts Inspector"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_inspector_website"
description: ""
---

# huaweicloud_codearts_inspector_website

Manages a CodeArts inspector website resource within HuaweiCloud.

## Example Usage

```hcl
variable "website_address" {}
variable "login_url" {}
variable "login_username" {}
variable "login_password" {}
variable "login_cookie" {}
variable "verify_url" {}

resource "huaweicloud_codearts_inspector_website" "test" {
  website_name    = "test-name"
  auth_type       = "free"
  website_address = var.website_address
  login_url       = var.login_url
  login_username  = var.login_username
  login_password  = var.login_password
  login_cookie    = var.login_cookie
  verify_url      = var.verify_url

  http_headers = {
    "test-key" = "test-value"
  }  
}
```

## Argument Reference

The following arguments are supported:

* `website_name` - (Required, String, ForceNew) Specifies the website name. The valid length is limited from `1` to `50`.
  Changing this parameter will create a new resource. Only Chinese characters, letters, digits, hyphens (-) and
  underscores (_) are allowed, and cannot start with a hyphen (-).

* `website_address` - (Required, String, ForceNew) Specifies the unique website address. The maximum length is `256`.
  Changing this parameter will create a new resource.
  The format should be `http(s)://example.com` or `http(s)://{public IPv4 address}:{PORT}`.

* `auth_type` - (Required, String, ForceNew) Specifies the authentication type. Changing this parameter will create a
  new resource. Valid values are:
  + **free**: Verification Free. Before using this authentication method, please confirm the following instructions.
  Please confirm that your account has completed real-name authentication and is not a restricted account. Please confirm
  that you have obtained the relevant legal rights to scan the scanned objects. Please confirm that your scanning behavior
  has legal and reasonable purposes and complies with applicable laws and regulations. Illegal scanning with this Service
  is not allowed. If there are any violations to the terms and conditions described here, Huawei is entitled to immediately
  terminate your use of this Service and you shall compensate us and any related third parties for any losses incurred therefrom.
  + **auto**: One-Click Verification. Before using this authentication method, please confirm that the server of the site
  to be detected is built in HuaweiCloud, and that the server is the asset of your current login account.

* `login_url` - (Optional, String) Specifies the login URL. The login address and domain name must be the same as the field
  `website_address`, for example: `http(s)://example.com/login`. The maximum length is `2048`.

* `login_username` - (Optional, String) Specifies the login username. The maximum length is `256`.

* `login_password` - (Optional, String) Specifies the login password. The maximum length is `50`.

-> Fields `login_url`, `login_username` and `login_password` must be set simultaneously.

* `login_cookie` - (Optional, String) Specifies the login cookie.

* `http_headers` - (Optional, Map) Specifies the custom HTTP request headers, the format is key/value pairs.

* `verify_url` - (Optional, String) Specifies the verify URL that can only be accessed after successful login.
  CodeArts inspector will use this URL to quickly determine whether your login information is valid.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `high` - The number of high-risk vulnerabilities.

* `middle` - The number of medium-risk vulnerabilities.

* `low` - The number of low-severity vulnerabilities.

* `hint` - The number of hint-risk vulnerabilities.

* `created_at` - The time to create website.

* `auth_status` - The auth status of website. Valid values are:
  + **unauth**: Unauthorized.
  + **auth**: Authorized.
  + **invalid**: Authentication file is invalid.
  + **manual**: Manual authentication.
  + **skip**: Authentication free.

## Import

The CodeArts inspector website can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_inspector_website.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `auth_type`, `login_password`, `login_cookie`,
`http_headers`. It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_inspector_website" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      auth_type,
      login_password,
      login_cookie,
      http_headers,
    ]
  }
}
```
