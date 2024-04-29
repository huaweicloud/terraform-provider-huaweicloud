---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_flink_template"
description: ""
---

# huaweicloud_dli_flink_template

Manages a DLI Flink template resource within HuaweiCloud.  

## Example Usage

```hcl
  variable "sql" {}
  
  resource "huaweicloud_dli_flink_template" "test" {
    name        = "demo"
    type        = "flink_sql_job"
    sql         = var.sql
    description = "This is a demo"

    tags = {
      foo = "bar"
      key = "value"
    }
  }
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The name of the flink template.

* `sql` - (Optional, String) The statement of the flink template.

* `description` - (Optional, String) The description of the flink template.

* `type` - (Optional, String, ForceNew) The type of the flink template.  
  Valid values are **flink_sql_job** and **flink_opensource_sql_job**.
  Defaults to **flink_sql_job**.

  Changing this parameter will create a new resource.

* `tags` - (Optional, Map, ForceNew) The key/value pairs to associate with the flink template.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The flink template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dli_flink_template.test 1231
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`tags`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also you can ignore changes as below.

```bash
resource "huaweicloud_dli_flink_template" "test" {
    ...

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}
```
