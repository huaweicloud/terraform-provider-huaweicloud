---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_account_enroll"
description: |-
  Manages an RGC account enroll resource within HuaweiCloud.
---

# huaweicloud_rgc_account_enroll

Manages an RGC account enroll resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "parent_organizational_unit_name" {}
variable "parent_organizational_unit_id" {}

resource "huaweicloud_rgc_account_enroll" "test"{
  managed_account_id            = "test"
  parent_organizational_unit_id = "test@terraform.com"
}
```

### Account with Blueprint

```hcl
variable "managed_account_id" {}
variable "parent_organizational_unit_id" {}
variable "blueprint_product_id" {}
variable "blueprint_product_version" {}

resource "huaweicloud_rgc_account" "test" {
  managed_account_id            = var.managed_account_id
  parent_organizational_unit_id = var.parent_organizational_unit_id

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

* `managed_account_id` - (Required, String, NonUpdatable) Specifies the ID of the account.

* `parent_organizational_unit_id` - (Required, String, NonUpdatable) Specifies organizational unit ID of enrolled account.

* `blueprint` - (Optional, List, NonUpdatable) Specifies the blueprint of the account.
  The [blueprint](#blueprint) structure is documented below.

<a name="blueprint"></a>
The `blueprint` block supports:

* `blueprint_product_id` - (Optional, String) Specifies the ID of the blueprint.

* `blueprint_product_version` - (Optional, String) Specifies the version of the blueprint.

* `parent_organizational_unit_name` - (Optional, String) Name of a registered parent OU.

* `variables` - (Optional, String) Specifies the variables of the blueprint.

* `is_blueprint_has_multi_account_resource` - (Optional, Bool) Specifies whether the blueprint has multi-account resources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `stage` - Indicates the state of enrolled account.

* `account_type` - The account enrolled type.

* `account_name` - The account enrolled name.

* `landing_zone_version` - Landing zone version of an enrolled account.

* `manage_account_id` - Management account ID.

* `owner` -  The account owner.

* `created_at` - The creation time.

* `updated_at` - The last update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 6 hours.

* `delete` - Default is 6 hours.

## Import

The RGC account enroll can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rgc_account_enroll.test <id>
```
