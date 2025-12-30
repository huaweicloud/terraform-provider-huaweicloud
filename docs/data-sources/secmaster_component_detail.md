---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_component_detail"
description: |-
  Use this data source to query a specific component detail.
---

# huaweicloud_secmaster_component_detail

Use this data source to query a specific component detail.

## Example Usage

```hcl
variable "workspace_id" {}
variable "component_id" {}

data "huaweicloud_secmaster_component_detail" "test" {
  workspace_id = var.workspace_id
  component_id = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `component_id` - (Required, String) Specifies the component ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `component_id_attr` - The component ID.

* `component_name` - The component name.

* `description` - The component description.

* `create_time` - The creation time.

* `history_version` - The history version.

* `maintainer` - The workflow trigger mode.

* `time_zone` - The time zone.

* `update_time` - The update time.

* `upgrade` - The upgrade status.
  The value can be **NEED_UPGRADE**, **UPGRADING**, **SUCCESS_UPGRADE** or **FAIL_UPGRADE**.

* `upgrade_fail_message` - The upgrade failed message.

* `version` - The SecMaster version.
