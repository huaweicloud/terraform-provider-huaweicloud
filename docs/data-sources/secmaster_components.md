---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_components"
description: |-
  Use this data source to get the list of SecMaster components within HuaweiCloud.
---

# huaweicloud_secmaster_components

Use this data source to get the list of SecMaster components within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_components" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `records` - The components list.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `component_id` - The component ID.

* `component_name` - The component name.

* `version` - The SecMaster version.

* `history_version` - The history version.

* `create_time` - The creation time (timestamp in milliseconds).

* `update_time` - The update time (timestamp in milliseconds).

* `time_zone` - The time zone.

* `upgrade` - The upgrade status. Valid values are **NEED_UPGRADE**, **UPGRADING**, **SUCCESS_UPGRADE**,
  **FAIL_UPGRADE**.

* `upgrade_fail_message` - The upgrade failure message.

* `maintainer` - The maintainer.

* `description` - The component description.
