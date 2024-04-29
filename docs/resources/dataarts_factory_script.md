---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_factory_script"
description: ""
---

# huaweicloud_dataarts_factory_script

Manages DataArts Factory script resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "name" {}
variable "connection_name" {}

resource "huaweicloud_dataarts_factory_script" "test" {
  workspace_id    = var.workspace_id
  name            = var.name
  type            = "DLISQL"
  content         = "#content"
  connection_name = var.connection_name
  queue_name      = "default"
  description     = "test"
  configuration   = {
    "spark.sql.files.maxRecordsPerFile" = "1"
    "dli.sql.job.timeout"               = "1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the script.
  Changing this creates a new script.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID which the script in.
  Changing this creates a new script.

* `name` - (Required, String, ForceNew) Specifies the script name. The name contains a maximum of 128 characters,
  including only letters, numbers, hyphens (-), and periods (.). The script name must be unique. Changing this creates
  a new script.

* `type` - (Required, String, ForceNew) Specifies the script type. The valid values are: **FlinkSQL**, **DLISQL**,
  **SparkSQL**, **HiveSQL**, **DWSSQL**, **RDSSQL**, **Shell**, **PRESTO**, **ClickHouseSQL**, **HetuEngineSQL**,
  **PYTHON**, **ImpalaSQL**. Changing this creates a new script.

* `content` - (Required, String) Specifies the script content. A maximum of 4 MB is supported.

* `connection_name` - (Required, String) Specifies the connection name of script.

* `directory` - (Optional, String) Specifies the directory of script.

* `database` - (Optional, String) Specifies the database of script.

* `queue_name` - (Optional, String) Specifies the queue name of script.

* `description` - (Optional, String) Specifies the description of script.

* `target_status` - (Optional, String) Specifies the target status of script.

* `configuration` - (Optional, Map) Specifies the configuration of script. Only valid key and value can take an effect.
 if put the invalid key and value in the `configuration` map, it may proceed with an empty `configuration` map.

* `approvers` - (Optional, List) Specifies the approvers of script.
The [approvers](#approvers) structure is documented below.

  <a name="approvers"></a>
  The `approvers` block supports:
* `approver_name` - (Required, String) Specifies the approver name of script.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The DataArts Factory script id. The format of the id is `<workspace_id>/<name>`.

* `auto_acquire_lock` - Whether the resource automatically obtain edit lock parameters.

* `created_by` - The person creating the script.

## Import

DataArts factory script can be imported using `<workspace_id>/<name>`, e.g.

```bash
terraform import huaweicloud_dataarts_factory_script.test b41a17b18a814b118730a867cecb9952/test
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `approvers`, `target_status`.

It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dataarts_factory_script" "test" {
  ...

  lifecycle {
    ignore_changes = [
      approvers, target_status,
    ]
  }
}
```
