---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_soc_attachment"
description: |-
  Use this data source to get the SecMaster soc attachment detail within HuaweiCloud.
---

# huaweicloud_secmaster_soc_attachment

Use this data source to get the SecMaster soc attachment detail within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "attach_id" {}

data "huaweicloud_secmaster_soc_attachment" "test" {
  workspace_id = var.workspace_id
  attach_id    = var.attach_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `attach_id` - (Required, String) Specifies the attachment ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The soc attachment detail.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The attachment ID.

* `file_name` - The attachment name.

* `file_folder` - The folder.

* `workspace_id` - The workspace ID.

* `creator_id` - The creator ID.

* `creator_name` - The creator name.

* `create_time` - The creation time.

* `update_time` - The update time.

* `modifier_id` - The modifier ID.

* `modifier_name` - The modifier name.

* `is_deleted` - Whether the attachment is deleted.

* `storage_type` - The storage type.

* `storage_url` - The storage URL.
