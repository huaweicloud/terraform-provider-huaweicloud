---
subcategory: "RGC"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_accounts"
description: |-
  Use this data source to list managed accounts in Resource Governance Center.
---

# huaweicloud_rgc_accounts

Use this data source to list managed accounts in Resource Governance Center.

## Example Usage

```hcl
variable control_id {}
data "huaweicloud_rgc_accounts" "test" {}

data "huaweicloud_rgc_accounts" "test_control_id" {
  control_id = var.control_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `control_id` - (Optional, String) The enabled control policy information.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `managed_accounts` - The managed account information.

  The [managed_accounts](#managed_accounts) structure is documented below.

<a name="managed_accounts"></a>
The `managed_accounts` block supports:

* `account_id` - The managed account ID.

* `account_type` - The managed account type.

* `blueprint_product_id` - The blueprint ID.

* `blueprint_product_version` - The blueprint version.

* `blueprint_status` - The blueprint deployment status, including failure, completion, and in progress.

* `created_at` - The time when the managed account is created in an OU.

* `identity_store_user_name` - The Identity Center user name.

* `is_blueprint_has_multi_account_resource` - Whether the blueprint contains multi-account resources.

* `landing_zone_version` - The Landing Zone version.

* `manage_account_id` - The management account ID.

* `message` - The error status description information.

* `owner` - The source of the managed account creation, including CUSTOM and RGC.

* `parent_organizational_unit_id` - The parent OU ID.

* `parent_organizational_unit_name` - The parent OU name.

* `state` - The managed account status.

* `updated_at` - The time when the managed account was last updated in an OU.

* `regions` - The region information.

  The [regions](#regions) structure is documented below.

<a name="regions"></a>
The `regions` block supports:

* `region` - The region name.

* `region_status` - The region status, which can be available or unavailable.
