---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_modules"
description: |-
  Use this data source to get the list of modules.
---

# huaweicloud_secmaster_modules

Use this data source to get the list of modules.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_modules" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `module_type` - (Optional, String) Specifies the module type.
  The value can be **section** or **tab**.

* `sort_key` - (Optional, String) Specifies sorting field.
 The value can be **create_time** (default value) or **update_time**.

* `sort_dir` - (Optional, String) Specifies sorting order.
  The valid values are as follows:
  + **ASC**
  + **DESC** (Defaults)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of modules.

  The [data](#modules_struct) structure is documented below.

<a name="modules_struct"></a>
The `data` block supports:

* `cloud_pack_id` - The subscription package ID.

* `cloud_pack_name` - The subscription package name.

* `cloud_pack_version` - The subscription package version.

* `create_time` - The creation time.

* `creator_id` - The creator ID.

* `description` - The module description.

* `en_description` - The module English description.

* `id` - The module ID.

* `module_json` - The module related information.

* `name` - The module name.

* `en_name` - The module English name.

* `project_id` - The project ID.

* `workspace_id` - The workspace ID.

* `update_time` - The update time.

* `thumbnail` - The module thumbnail.

* `module_type` - The module type.

* `tag` - The module tag.

* `is_built_in` - Whether the module is a system module.

* `data_query` - The data query method.

* `boa_version` - The BOA version.

* `version` - The SecMaster version.
