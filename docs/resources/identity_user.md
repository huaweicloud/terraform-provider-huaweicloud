---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud\_identity\_user

Manages a User resource within HuaweiCloud IAM service.
This is an alternative to `huaweicloud_identity_user_v3`

Note: You _must_ have admin privileges in your HuaweiCloud cloud to use
this resource.

## Example Usage

```hcl
resource "huaweicloud_identity_user" "user_1" {
  name        = "user_1"
  description = "A user"
  password    = "password123!"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the user. The user name consists of 5 to 32
     characters. It can contain only uppercase letters, lowercase letters, 
     digits, spaces, and special characters (-_) and cannot start with a digit.

* `description` - (Optional, String) A description of the user.

* `default_project_id` - (Optional, String) The default project this user belongs to.

* `domain_id` - (Optional, String) The domain this user belongs to.

* `enabled` - (Optional, Bool) Whether the user is enabled or disabled. Valid
    values are `true` and `false`.

* `password` - (Optional, String) The password for the user. It must contain at least 
     two of the following character types: uppercase letters, lowercase letters, 
     digits, and special characters.

## Attributes Reference

The following attributes are exported:

## Import

Users can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_identity_user_v3.user_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
