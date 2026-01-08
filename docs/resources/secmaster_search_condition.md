---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_search_condition"
description: |-
  Manages a search condition resource within HuaweiCloud.
---

# huaweicloud_secmaster_search_condition

Manages a search condition resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "dataspace_id" {}
variable "pipe_id" {}
variable "condition_name" {}
variable "query" {}

resource "huaweicloud_secmaster_search_condition" "test" {
  workspace_id   = var.workspace_id
  dataspace_id   = var.dataspace_id
  pipe_id        = var.pipe_id
  condition_name = var.condition_name
  query          = var.query
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the dataspace belongs.

* `dataspace_id` - (Required, String, NonUpdatable) Specifies the ID of the dataspace.

* `pipe_id` - (Required, String, NonUpdatable) Specifies the ID of the pipe.

* `condition_name` - (Required, String) Specifies the name of the search condition.

* `query` - (Required, String) Specifies the query statement.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The workflow can be imported using the `workspace_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_search_condition.test <workspace_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `dataspace_id`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_secmaster_search_condition" "test" {
  ...

  lifecycle {
    ignore_changes = [
      dataspace_id,
    ]
  }
}
```
