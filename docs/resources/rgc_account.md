---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_account"
description: |-
  Manages an RGC account resource within HuaweiCloud.
---

# huaweicloud_rgc_account

Manages an RGC account resource within HuaweiCloud.

-> **NOTE:** Deleting RGC account is not support. If you destroy a resource of RGC account,
the RGC account is only closed and removed from the state, but it remains in the cloud.

## Example Usage

### Basic Usage

```hcl
variable "parent_organizational_unit_name" {}
variable "parent_organizational_unit_id" {}

resource "huaweicloud_rgc_account" "test"{
  name                            = "test"
  email                           = "test@terraform.com"
  parent_organizational_unit_name = var.parent_organizational_unit_name
  parent_organizational_unit_id   = var.parent_organizational_unit_id
  identity_store_user_name        = "test"
  identity_store_email            = "test@terraform.com"
}
```

### Account with Blueprint

```hcl
variable "parent_organizational_unit_name" {}
variable "parent_organizational_unit_id" {}
variable "blueprint_product_id" {}
variable "blueprint_product_version" {}

resource "huaweicloud_rgc_account" "test"{
  name                            = "test"
  email                           = "test@terraform.com"
  parent_organizational_unit_name = var.parent_organizational_unit_name
  parent_organizational_unit_id   = var.parent_organizational_unit_id
  identity_store_user_name        = "test"
  identity_store_email            = "test@terraform.com"

  blueprint {
    blueprint_product_id                    = var.blueprint_product_id
    blueprint_product_version               = var.blueprint_product_version
    is_blueprint_has_multi_account_resource = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the account.

* `email` - (Required, String, NonUpdatable) Specifies the email address of the account.

* `phone` - (Optional, String, NonUpdatable) Specifies the mobile number of the account.

* `parent_organizational_unit_name` - (Required, String, NonUpdatable) Specifies the name of the parent organizational unit.

* `parent_organizational_unit_id` - (Required, String ,NonUpdatable) Specifies the ID of the parent organizational unit.

* `identity_store_user_name` - (Required, String, NonUpdatable) Specifies the name of the account in identity center.

* `identity_store_email` - (Required, String, NonUpdatable) Specifies the email address of the account in identity center.

* `blueprint` - (Optional, List, NonUpdatable) Specifies the blueprint of the account.
  The [blueprint](#blueprint) structure is documented below.

<a name="blueprint"></a>
The `blueprint` block supports:

* `blueprint_product_id` - (Optional, String) Specifies the ID of the blueprint.

* `blueprint_product_version` - (Optional, String) Specifies the version of the blueprint.

* `variables` - (Optional, String) Specifies the variables of the blueprint.

* `is_blueprint_has_multi_account_resource` - (Optional, Bool) Specifies whether the blueprint has multi-account resources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - Indicates the uniform resource name of the account.

* `joined_at` - Indicates the time when the account was created.

* `joined_method` - Indicates how an account joined an organization.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 15 minutes.
* `delete` - Default is 15 minutes.

## Import

The Organizations account can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rgc_account.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include `email`, `phone`, `parent_organizational_unit_name`, `parent_organizational_unit_id`,
`identity_store_user_name`, `identity_store_email` and `blueprint`.
It is generally recommended running `terraform plan` after importing an account.
You can then decide if changes should be applied to the account, or the resource definition should be updated to
align with the account. Also you can ignore changes as below.

```hcl
resource "huaweicloud_rgc_account" "test" {
  ...

  lifecycle {
    ignore_changes = [
      email, phone, parent_organizational_unit_name, parent_organizational_unit_id, identity_store_user_name,
      identity_store_email, blueprint
    ]
  }
}
```
