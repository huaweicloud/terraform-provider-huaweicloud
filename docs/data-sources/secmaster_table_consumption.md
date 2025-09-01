---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_table_consumption"
description: |-
  Use this data source to get the data consumption configuration of a specified table from HuaweiCloud SecMaster.
---

# huaweicloud_secmaster_table_consumption

Use this data source to get the data consumption configuration of a specified table from HuaweiCloud SecMaster.

## Example Usage

```hcl
variable "workspace_id" {}
variable "table_id" {}

data "huaweicloud_secmaster_table_consumption" "example" {
  workspace_id = var.workspace_id
  table_id     = var.table_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `table_id` - (Required, String) Specifies the ID of the table to get consumption configuration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `table_name` - The name of the table.

* `dataspace_id` - The data space ID.

* `pipe_id` - The pipeline ID.

* `pipe_name` - The name of the pipeline.

* `status` - The status of the consumption configuration. Possible values are **ENABLED** or **DISABLED**.

* `type` - The network type. Possible values are **INTERNET** or **INTRANET**.

* `access_point` - The access point domain information (format: `{dataspace}.{endpoint}`).

* `subscirption` - The subscription name.
