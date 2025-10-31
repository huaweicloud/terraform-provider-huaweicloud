---
subcategory: "RGC"
layout: "huaweicloud"
page_title: "huaweicloud_rgc_organizational_unit_accounts"
description: |
  Use this data source to list managed accounts for a managed organizational unit in Resource Governance Center.
---

# huaweicloud_rgc_organizational_unit_accounts

Use this data source to list managed accounts for a managed organizational unit in Resource Governance Center.

## Example Usage

```hcl
variable "managed_organizational_unit_id" {}

data "huaweicloud_rgc_organizational_unit_accounts" "test" {
  managed_organizational_unit_id = var.managed_organizational_unit_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `managed_organizational_unit_id - (Required, String) The ID of the managed organizational unit for which
  to retrieve managed accounts.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `managed_accounts` - Information about the managed accounts for the managed organizational unit.

 The [managed_accounts](#managed_accounts) structure is documented below.

<a name="managed_accounts"></a>
The `managed_accounts` block supports:

* `landing_zone_version` - The version of the landing zone.

* `manage_account_id` - The ID of the managing account.

* `account_id` - The ID of the managed account.

* `account_name` - The name of the managed account.

* `account_type` - The type of the managed account.

* `owner` - The owner of the managed account.

* `state` - The state of the managed account.

* `message` - A message associated with the managed account.

* `parent_organizational_unit_id` - The ID of the parent organizational unit.

* `parent_organizational_unit_name` - The name of the parent organizational unit.

* `identity_store_user_name` - The username from the identity store.

* `blueprint_product_id` - The ID of the blueprint product.

* `blueprint_product_version` - The version of the blueprint product.

* `blueprint_status` - The status of the blueprint.

* `is_blueprint_has_multi_account_resource` - Indicates whether the blueprint has multi-account resources.

* `created_at` - The timestamp when the account was created.

* `updated_at` - The timestamp when the account was last updated.

* `regions` - A list of regions associated with the account.

The [regions](#regions) structure is documented below.

<a name="regions"></a>
The `regions` block supports:

* `region` - The name of the region.

* `region_status` - The status of the region.
