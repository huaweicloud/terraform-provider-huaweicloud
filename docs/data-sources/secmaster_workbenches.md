---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workbenches"
description: |-
  Use this data source to get the list of SecMaster workbenches within HuaweiCloud.
---

# huaweicloud_secmaster_workbenches

Use this data source to get the list of SecMaster workbenches within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_workbenches" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `name` - (Optional, String) Specifies the workbench name used for filtering.

* `status` - (Optional, String) Specifies the workbench status used for filtering.
  The valid values are as follows:
  + **publish**: Published.
  + **unpublish**: Unpublished.

* `type` - (Optional, String) Specifies the workbench type used for filtering.
  The valid values are as follows:
  + **scenario**: Scenario workbench.
  + **defense**: Defense workbench.

* `creator_type` - (Optional, String) Specifies the creator type used for filtering.
  The valid values are as follows:
  + **system**: System created.
  + **mine**: Created by me.
  + **others**: Created by others.

* `global_search_text` - (Optional, String) Specifies the global search keyword for filtering.

* `tags` - (Optional, String) Specifies the workbench tag used for filtering.

* `from_date` - (Optional, String) Specifies the start time of the creation time range for filtering.

* `to_date` - (Optional, String) Specifies the end time of the creation time range for filtering.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The workbenches list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The workbench ID.

* `name` - The workbench name.

* `type` - The workbench type.
  The valid values are as follows:
  + **DEFENSE**: Defense workbench.
  + **SCENARIO**: Scenario workbench.

* `url` - The workbench homepage address.

* `url_openwith_type` - The opening method of the workbench homepage.
  The valid values are as follows:
  + **NEW_PAGE**: Open in a new page.
  + **CURRENT_PAGE**: Open in the current page.

* `tags` - The workbench tags, separated by commas.

* `description` - The workbench description.

* `icon` - The workbench icon.

* `basic_properties` - The basic properties of the workbench.

* `domain_id` - The account ID.

* `region_id` - The region ID.

* `workspace_id` - The workspace ID.

* `create_time` - The creation time of the workbench.

* `update_time` - The update time of the workbench.

* `creator_id` - The creator ID.

* `creator_name` - The creator name.

* `modifier_id` - The modifier ID.

* `modifier_name` - The modifier name.

* `is_deleted` - Whether the workbench is deleted.

* `is_favorite` - Whether the workbench is favorited.

* `status` - The workbench publish status.
  The valid values are as follows:
  + **PUBLISH**: Published.
  + **UNPUBLISH**: Unpublished.
