---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_data_object_relations"
description: |-
  Manages a SecMaster data object relations resource within HuaweiCloud.
---

# huaweicloud_secmaster_data_object_relations

Manages a SecMaster data object relations resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "data_object_id" {}
variable "related_data_object_ids" {}

resource "huaweicloud_secmaster_data_object_relations" "test" {
  workspace_id            = var.workspace_id
  data_class              = "alerts"
  data_object_id          = var.data_object_id
  related_data_class      = "incidents"
  related_data_object_ids = var.related_data_object_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the data object belongs.

* `data_class` - (Required, String, NonUpdatable) Specifies the data class to which the data object belongs.
  The value can be **incidents**, **alerts** and **indicators**.

* `data_object_id` - (Required, String, NonUpdatable) Specifies the ID of the data object.

* `related_data_class` - (Required, String, NonUpdatable) Specifies the data class to which the related data object belongs.
  The value can be **incidents**, **alerts** and **indicators**.

* `related_data_object_ids` - (Required, List) Specifies the IDs of related the data objects.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The SecMaster data object relations can be imported using `workspace_id`, `data_class`, `data_object_id` and
`related_data_class` separated by slashs, e.g.

```bash
$ terraform import huaweicloud_secmaster_data_object_relations.test <workspace_id>/<data_class>/<data_object_id>/<related_data_class>
```
