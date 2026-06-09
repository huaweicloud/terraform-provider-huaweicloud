---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_compliance_packages"
description: |-
  Use this data source to get the list of SecMaster compliance packages within HuaweiCloud.
---

# huaweicloud_secmaster_compliance_packages

Use this data source to get the list of SecMaster compliance packages within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_compliance_packages" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `name` - (Optional, String) Specifies the compliance package name.

* `description` - (Optional, String) Specifies the compliance package description.

* `type` - (Optional, Int) Specifies the compliance package type.
  The value can be **0** (built-in) or **1** (customized).

* `state` - (Optional, Int) Specifies the compliance package state.
  The value can be **0** (disabled) or **1** (enabled).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `builtin_compliance_num` - The number of built-in compliance packages.

* `customized_compliance_num` - The number of customized compliance packages.

* `disabled_compliance_num` - The number of disabled compliance packages.

* `enabled_compliance_num` - The number of enabled compliance packages.

* `compliance_packages` - The list of compliance packages.

  The [compliance_packages](#compliance_packages_struct) structure is documented below.

<a name="compliance_packages_struct"></a>
The `compliance_packages` block supports:

* `uuid` - The UUID of the compliance package.

* `name` - The name of the compliance package.

* `version` - The version of the compliance package.

* `owner` - The owner of the compliance package.

* `description` - The description of the compliance package.

* `classify` - The classification of the compliance package.

* `areas` - The applicable areas of the compliance package.

* `region` - The applicable region of the compliance package.

* `state` - The state of the compliance package.

* `type` - The type of the compliance package.

* `check_items_num` - The number of check items in the compliance package.

* `has_auto_check_items` - Whether the compliance package contains auto check items.

* `spec_catalog_vo_list` - The catalog list of the compliance package.

  The [spec_catalog_vo_list](#spec_catalog_vo_list_struct) structure is documented below.

<a name="spec_catalog_vo_list_struct"></a>
The `spec_catalog_vo_list` block supports:

* `uuid` - The UUID of the catalog.

* `serial_number` - The serial number of the catalog.

* `level_number` - The level number of the catalog.

* `root` - The root UUID of the catalog.

* `parent` - The parent UUID of the catalog.

* `is_leaf` - Whether the catalog is a leaf node.

* `check_items` - The check items of the catalog.

  The [check_items](#check_items_struct) structure is documented below.

<a name="check_items_struct"></a>
The `check_items` block supports:

* `uuid` - The UUID of the check item.

* `name` - The name of the check item.
