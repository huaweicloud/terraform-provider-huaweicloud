---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_snapshot_copy"
description: |- 
  Use this resource to copy an automated snapshot to a manual snapshot within HuaweiCloud.
---

# huaweicloud_dws_snapshot_copy

Use this resource to copy an automated snapshot to a manual snapshot within HuaweiCloud.

-> 1. An automated snapshot can only correspond to one `huaweicloud_dws_snapshot_copy` resource.
   <br>2. Deleting this resource will delete the corresponding copied snapshot.

## Example Usage

```hcl
variable "automated_snapshot_id" {}
variable "snapshot_name" {}
variable "description" {}

resource "huaweicloud_dws_snapshot_copy" "test" {
  snapshot_id = var.automated_snapshot_id
  name        = var.snapshot_name
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `snapshot_id` - (Required, String, ForceNew) Specifies the ID of the automated snapshot to be copied.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the copy snapshot.
  Changing this creates a new resource.  
  The valid length is limited from `4` to `64`, only English letters, digits, underscores (_) and hyphens (-) are allowed,
  and must start with a letter.  
  The snapshot name must be unique.

* `description` - (Optional, String, ForceNew) Specifies the description of the copy snapshot.
  Changing this creates a new resource.  
  The maximum length is limited to `256` characters, and special characters (`!<>'=&"`) are not allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the ID of the copied snapshot.

## Import

The resource can be imported using the related `snapshot_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dws_snapshot_copy.test <snapshot_id>/<id>
```
