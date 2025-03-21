---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_search_criteria"
description: |-
  Manages an LTS search criteria resource within HuaweiCloud.
---

# huaweicloud_lts_search_criteria

Manages an LTS search criteria resource within HuaweiCloud.

## Example Usage

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}

resource "huaweicloud_lts_search_criteria" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id

  criteria = "content:test"
  name     = "search_criteria_test"
  type     = "ORIGINALLOG"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `log_group_id` - (Required, String, ForceNew) Specifies the ID of a log group. Changing this parameter will create
  a new resource.

* `log_stream_id` - (Required, String, ForceNew) Specifies the ID of a log stream. Changing this parameter will create
  a new resource.

* `criteria` - (Required, String, ForceNew) Specifies the content of search criteria. Changing this parameter will create
  a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the search criteria. The name can only contain English
  letters, numbers, Chinese characters, hyphens, underscores, and periods. It cannot start with a period or underscore
  or end with a period.Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the type of the search criteria. Available types are
  **ORIGINALLOG** (for raw logs) and **VISUALIZATION** (for visualized logs). Changing this parameter will create a new
  resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The search criteria can be imported using the group ID, stream ID, and resource ID separated by the slashes, e.g.

```bash
$ terraform import huaweicloud_lts_search_criteria.test <log_group_id>/<log_stream_id>/<id>
```
