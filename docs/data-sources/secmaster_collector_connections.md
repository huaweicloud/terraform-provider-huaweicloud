---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_connections"
description: |-
  Use this data source to get the list of collector connections.
---

# huaweicloud_secmaster_collector_connections

Use this data source to get the list of collector connections.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_collector_connections" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID to which the collector connections belong.

* `connection_type` - (Optional, String) Specifies the connection type. Valid values are:
  + **FILTER**
  + **INPUT**
  + **OUTPUT**

* `title` - (Optional, String) Specifies the title of the collector connection.

* `description` - (Optional, String) Specifies the description of the collector connection.

* `sort_key` - (Optional, String) Specifies the sort key of the collector connection.

* `sort_dir` - (Optional, String) Specifies the sort direction of the collector connection.
  Valid values are **asc** and **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of collector connection details.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `channel_refer_count` - The channel reference count of the collector connection.

* `connection_id` - The ID of the collector connection.

* `connection_type` - The connection type of the collector connection.

* `description` - The description of the collector connection.

* `info` - The info of the collector connection.

* `module_id` - The module ID of the collector connection.

* `template_title` - The template title of the collector connection.

* `title` - The title of the collector connection.
