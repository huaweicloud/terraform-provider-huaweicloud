---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_pipe_index"
description: |-
  Manages a pipe index resource within HuaweiCloud.
---

# huaweicloud_secmaster_pipe_index

Manages a pipe index resource within HuaweiCloud.

-> System pipes do not support configuration.

## Example Usage

```hcl
variable "workspace_id" {}
variable "pipe_id" {}

resource "huaweicloud_secmaster_pipe_index" "test" {
  workspace_id    = var.workspace_id
  pipe_id         = var.pipe_id
  mapping         = jsonencode({
    field1 = {
      type = "text"
    }
    field2 = {
      type = "keyword"
    }
  })
  timestamp_field = "timestamp"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the pipe index
  belongs.

* `pipe_id` - (Required, String, NonUpdatable) Specifies the pipe ID of the pipe index.

* `mapping` - (Required, String) Specifies the index mapping information in JSON format.
  This is a JSON string representing a map structure where each key is a field name and each value is a
  [KeyIndex](#pipe_index_key_index) object describing the field's index configuration.

* `timestamp_field` - (Required, String) Specifies the timestamp field name.

<a name="pipe_index_key_index"></a>
The `KeyIndex` object supports:

* `type` - (Optional, String) Specifies the field type. The value can be **text** (full-text index field),
  **keyword** (exact match), **long** (long integer), **integer** (integer), **double** (double-precision floating-point),
  **float** (single-precision floating-point), or **date** (date type).

* `is_chinese_exist` - (Optional, Bool) Specifies whether the field contains Chinese characters.

* `properties` - (Optional, Map) Specifies the nested structure of the field, which is a map of field names
  to [KeyIndex](#pipe_index_key_index) objects. This field is used to define nested object structures.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the value of `pipe_id`.

## Import

The pipe index can be imported using the `workspace_id` and their `pipe_id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_pipe_index.test <workspace_id>/<pipe_id>
```
