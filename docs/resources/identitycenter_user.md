---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_user"
description: ""
---

# huaweicloud_identitycenter_user

Manages an Identity Center user resource within HuaweiCloud.

## Example Usage

```hcl
variable "identity_store_id" {}

resource "huaweicloud_identitycenter_user" "test"{
  identity_store_id = var.identity_store_id
  user_name         = "test_user"
  password_mode     = "OTP"
  display_name      = "test_display_name"
  family_name       = "test_family_name"
  given_name        = "test_given_name"
  email             = "email@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `identity_store_id` - (Required, String, ForceNew) Specifies the ID of the identity store.

  Changing this parameter will create a new resource.

* `user_name` - (Required, String, ForceNew) Specifies the username of the user.

  Changing this parameter will create a new resource.

* `password_mode` - (Required, String, ForceNew) Specifies the initialized password mode. Value options:
  + **OTP**: Generate random One-time password.
  + **EMAIL**: Sending an email to user which include password setting instructions.

  Changing this parameter will create a new resource.

* `family_name` - (Required, String) Specifies the family name of the user.

* `given_name` - (Required, String) Specifies the given name of the user.

* `display_name` - (Required, String) Specifies the display name of the user.

* `email` - (Required, String) Specifies the email of the user.

* `phone_number` - (Optional, String) Specifies the phone number of the user.

* `user_type` - (Optional, String) Specifies the type of the user.

* `title` - (Optional, String) Specifies the title of the user.

* `addresses` - (Optional, List) Specifies the addresses information of the user.
  The [addresses](#addresses_struct) structure is documented below.

* `enterprise` - (Optional, List) Specifies the enterprise information of the user.
  The [enterprise](#enterprise_struct) structure is documented below.

* `enabled` - (Optional, Bool) Whether the user is enabled. Value options:
  + **true**
  + **false**

<a name="addresses_struct"></a>
The `addresses` block supports:

* `country` - (Optional, String) Specifies the country of the user.

* `formatted` - (Optional, String) Specifies a string containing a formatted version of the address to be displayed.

* `locality` - (Optional, String) Specifies the locality of the user.

* `postal_code` - (Optional, String) Specifies the postal code of the user.

* `region` - (Optional, String) Specifies the region of the user.

* `street_address` - (Optional, String) Specifies the street address of the user.

<a name="enterprise_struct"></a>
The `enterprise` block supports:

* `cost_center` - (Optional, String) Specifies the cost center of the enterprise.

* `department` - (Optional, String) Specifies the department of the enterprise.

* `division` - (Optional, String) Specifies the division of the enterprise.

* `employee_number` - (Optional, String) Specifies the employee number of the enterprise.

* `organization` - (Optional, String) Specifies the organization of the enterprise.

* `manager` - (Optional, String) Specifies the manager of the enterprise.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time of the user.

* `created_by` - The creator of the user.

* `updated_at` - The update time of the user.

* `updated_by` - The updater of the user.

* `email_verified` - Whether the email is verified.

## Import

The IdentityCenter user can be imported using the `identity_store_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_user.test <identity_store_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password_mode`. It is generally
recommended running `terraform plan` after importing an IdentityCenter user. You can then decide if changes should be
applied to the IdentityCenter user, or the resource definition should be updated to align with the instance. Also, you
can ignore changes as below.

```hcl
resource "huaweicloud_identitycenter_user" "user" {
  ...

  lifecycle {
    ignore_changes = [
      password_mode,
    ]
  }
}
```
