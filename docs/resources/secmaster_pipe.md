---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_pipe"
description: |-
  Manages a data pipe resource within HuaweiCloud SecMaster.
---

# huaweicloud_secmaster_pipe

Manages a data pipe resource within HuaweiCloud SecMaster.

## Example Usage

```hcl
variable "workspace_id" {}
variable "dataspace_id" {}

resource "huaweicloud_secmaster_pipe" "test" {
  workspace_id    = var.workspace_id
  dataspace_id    = var.dataspace_id
  pipe_name       = "test01"
  shards          = 2
  storage_period  = 20
  description     = "test description"
  timestamp_field = "timestamp"

  mapping = jsonencode({
    id = {
      is_chinese_exist = true
      properties       = {}
      type             = "text"
    }
    name = {
      is_chinese_exist = false
      properties       = {}
      type             = "text"
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the pipe belongs.

* `dataspace_id` - (Required, String, NonUpdatable) Specifies the data space ID.

* `pipe_name` - (Required, String, NonUpdatable) Specifies the name of the data pipe.
  The name must start with a letter and contain only lowercase letters, digits, and asterisks (*).
  Asterisks cannot appear at the end or consecutively.
  The name cannot start with system reserved prefixes: isap_, csb_, secmaster_, sec_, s_sec_, i_sec_, l_sec_, security_.

* `shards` - (Required, Int) Specifies the number of partitions for the data pipe.
  The value ranges from `1` to `64`.

* `storage_period` - (Required, Int) Specifies the data retention period in days.
  The default value is `30`, and the value ranges from `7` to `180`.

* `description` - (Optional, String) Specifies the description of the data pipe.

* `mapping` - (Optional, String, NonUpdatable) Specifies the index field mapping in JSON format.
  Each key object carries information about one field.
  Each field must contain `type` (string), `is_chinese_exist` (boolean), and `properties` (object).

  The `type` parameter valid values are as follows:
  + **text**: Full-text index field for text search.
  + **keyword**: Keyword type for exact matching.
  + **long**: Long integer type.
  + **integer**: Integer type.
  + **double**: Double-precision floating-point number.
  + **float**: Single-precision floating-point number.
  + **date**: Date type.

  Example: `{"field1": {"type": "text", "is_chinese_exist": true, "properties": {}}, "field2": {...}}`

* `timestamp_field` - (Optional, String, NonUpdatable) Specifies the timestamp field.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the pipe ID.

* `pipe_alias` - The alias of the pipe.

* `pipe_type` - The type of the pipe. The value can be **system-defined** or **user-defined**.

* `dataspace_name` - The name of the data space.

* `category` - The resource type.

* `owner_type` - The owner type.

* `process_status` - The processing status.

* `create_by` - The creator.

* `create_time` - The creation time.

* `update_by` - The updater.

* `update_time` - The update time.

* `domain_id` - The domain ID.

* `project_id` - The project ID.

## Import

The data pipe can be imported using the `workspace_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_pipe.test <workspace_id>/<pipe_id>
```
