---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_component_alliances"
description: |-
  Use this data source to get the list of SecMaster component alliances within HuaweiCloud.
---

# huaweicloud_secmaster_component_alliances

Use this data source to get the list of SecMaster component alliances within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_component_alliances" "test" {
  workspace_id = var.workspace_id
  is_built_in  = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `is_built_in` - (Required, Bool) Specifies whether to query only built-in component alliances.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of component alliances.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `alliance_code` - The built-in code of the component alliance.

* `alliance_description` - The description of the component alliance.

* `alliance_name` - The name of the component alliance.

* `alliance_type` - The type of the component alliance.
  The value can be **system** (built-in) or **custom** (customized).

* `logo` - The logo URL of the component alliance.

* `id` - The unique ID of the component alliance.

* `create_time` - The creation time of the component alliance.

* `creator_name` - The creator name of the component alliance.

* `update_time` - The update time of the component alliance.
