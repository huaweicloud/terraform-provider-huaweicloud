---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_layouts"
description: |-
  Use this data source to get the list of SecMaster layouts within HuaweiCloud.
---

# huaweicloud_secmaster_layouts

Use this data source to get the list of SecMaster layouts within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_layouts" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `name` - (Optional, String) Specifies the name of the layout to query.

* `used_by` - (Optional, String) Specifies the layout type.

* `binding_code` - (Optional, String) Specifies the binding code of the layout.

* `is_built_in` - (Optional, Bool) Specifies whether it is built into the system.

* `is_template` - (Optional, Bool) Specifies whether it is a template.

* `is_default` - (Optional, Bool) Specifies whether to use the default layout.

* `layout_type` - (Optional, String) Specifies the type of the layout.

* `sort_key` - (Optional, String) Specifies the key for sorting the results.

* `sort_dir` - (Optional, String) Specifies the direction for sorting the results.

* `search_txt` - (Optional, String) Specifies the search text to filter layouts.

* `from_date` - (Optional, String) Specifies the start date for filtering layouts.

* `to_date` - (Optional, String) Specifies the end date for filtering layouts.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The list of layouts.
  The [data](#secmaster_layouts_data) structure is documented below.

<a name="secmaster_layouts_data"></a>
The `data` block supports:

* `cloud_pack_id` - The cloud pack ID.

* `cloud_pack_name` - The cloud pack name.

* `cloud_pack_version` - The cloud pack version.

* `is_built_in` - Whether the layout is built-in.

* `is_template` - Whether the layout is a template.

* `create_time` - The creation time of the layout.

* `creator_id` - The ID of the user who created the layout.

* `parent_id` - The parent layout ID.

* `creator_name` - The name of the user who created the layout.

* `description` - The description of the layout.

* `en_description` - The English description of the layout.

* `id` - The ID of the layout.

* `name` - The name of the layout.

* `en_name` - The English name of the layout.

* `layout_json` - The JSON configuration of the layout.

* `project_id` - The project ID.

* `update_time` - The last update time of the layout.

* `workspace_id` - The workspace ID.

* `region_id` - The region ID.

* `domain_id` - The domain ID.

* `thumbnail` - The template thumbnail, this field has a value when the layout is a template.

* `used_by` - The usage scenario of the layout. Valid value are:
  + **DATACLASS**
  + **AOP_WORKFLOW**
  + **SECURITY_REPORT**
  + **DASHBOARD**

* `layout_cfg` - The front end binds an icon based on this value.

* `layout_type` - The type of the layout.
  This field is empty when `used_by` is **SECURITY_REPORT** or **DASHBOARD**.

* `binding_id` - The data class ID or process ID.
  This field is empty when `used_by` is **SECURITY_REPORT** or **DASHBOARD**.

* `binding_name` - The data class name or process name.
  This field is empty when `used_by` is **SECURITY_REPORT** or **DASHBOARD**.

* `binding_code` - The english names of data categories or processes.
  This field is empty when `used_by` is **SECURITY_REPORT** or **DASHBOARD**.

* `fields_sum` - The number of fields in the layout.
  This field is empty when `used_by` is **SECURITY_REPORT** or **DASHBOARD**.

* `wizards_sum` - The number of wizards in the layout.
  This field is empty when `used_by` is **SECURITY_REPORT** or **DASHBOARD**.

* `sections_sum` - The total number of system blocks.

* `modules_sum` - The total number of system modules.

* `tabs_sum` - The total number of custom indicators.

* `version` - The version of the SecMaster.

* `boa_version` - The BOA base version.
