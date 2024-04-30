---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_studio_data_connections"
description: |-
  Use this data source to get the list of the DataArts Studio data connections.
---

# huaweicloud_dataarts_studio_data_connections

Use this data source to get the list of the DataArts Studio data connections.

## Example Usage

```hcl
variable workspace_id {}

data "huaweicloud_dataarts_studio_data_connections" "test" {
  workspace_id  = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the data connection belongs.

* `connection_id` - (Optional, String) Specifies the ID of the data connection.

* `name` - (Optional, String) Specifies the name of the data connection.
  Supports fuzzy search.

* `type` - (Optional, String) Specifies the type of the data connection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - The list of the data connections.
  The [connections](#data_connections) structure is documented below.

<a name="data_connections"></a>
The `connections` block supports:

* `id` - The ID of the data connection.

* `name` - The name of the data connection.

* `type` - The type of the data connection.

* `agent_id` - The agent ID corresponding to the data connection.

* `qualified_name` - The qualified name of the data connection.

* `created_by` - The creator of the data connection.

* `created_at` - The creation time of the data connection.
