---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_channel_groups"
description: |-
  Use this data source to get the list of collector channel groups.
---

# huaweicloud_secmaster_collector_channel_groups

Use this data source to get the list of collector channel groups.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_collector_channel_groups" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID to which the collector channel groups belong.

* `name` - (Optional, String) Specifies the name of the collector channel group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of report details.

  The [groups](#groups_struct) structure is documented below.

<a name="groups_struct"></a>
The `groups` block supports:

* `group_id` - The ID of the collector channel group.

* `name` - The name of the collector channel group.
