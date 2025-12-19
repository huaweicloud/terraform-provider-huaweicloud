---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_dataspace"
description: |-
  Manages a SecMaster dataspace resource within HuaweiCloud.
---

# huaweicloud_secmaster_dataspace

Manages a SecMaster dataspace resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "dataspace_name" {}
variable "description" {}

resource "huaweicloud_secmaster_dataspace" "test" {
  workspace_id   = var.workspace_id
  dataspace_name = var.workflow_id
  description    = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Sepcifies the ID of the workspace to which the dataspace belongs.

* `dataspace_name` - (Required, String, NonUpdatable) Sepcifies the name of the dataspace.
  The name can only contain English letters, digits and hyphens (-), and cannot start or end with a hyphens (-),
  nor can they appear consecutively.
  The name valid length is limited from `5` to `63`.

* `description` - (Required, String) Sepcifies the description of the dataspace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `dataspace_type` - The dataspace type.
  The valid values are as follows:
  + **system-defined**
  + **user-defined**

* `domain_id` - The account ID.

* `project_id` - The project ID.

* `create_by` - The dataspace creator.

* `update_by` - The dataspace updater.

* `create_time` - The dataspace creation time.

* `update_time` - The dataspace update time.

## Import

The dataspace can be imported using the `workspace_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_dataspace.test <workspace_id>/<id>
```
