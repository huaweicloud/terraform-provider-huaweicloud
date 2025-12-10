---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_notes"
description: |-
  Use this data source to get the list of notes.
---

# huaweicloud_secmaster_notes

Use this data source to get the list of notes.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_notes" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID to which the notes belong.

* `sort_by` - (Optional, String) Specifies the field to sort by.

* `order` - (Optional, String) Specifies the order to sort by.

* `from_date` - (Optional, String) Specifies the start time to sort by.

* `to_date` - (Optional, String) Specifies the end time to sort by.

* `war_room_id` - (Optional, String) Specifies the war room ID to which the notes belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of notes.
  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `create_time` - The create time of the note.

* `update_time` - The update time of the note.

* `data` - The comment content body.
  The [data](#data_sub_struct) structure is documented below.

* `id` - The ID of the note.

* `is_deleted` - Whether the note is deleted.

* `marked_note` - Whether the note is marked.

* `note_type` - The type of the note.

* `project_id` - The project ID of the note.

* `type` - The type of the note.

* `user` - The user information.
  The [user](#user_struct) structure is documented below.

* `war_room_id` - The war room ID of the note.

* `workspace_id` - The workspace ID of the note.

* `content` - The content detail in JSON format string.

<a name="data_sub_struct"></a>
The `data` block supports:

* `content` - The content of the note.

<a name="user_struct"></a>
The `user` block supports:

* `id` - The ID of the user.

* `name` - The name of the user.
