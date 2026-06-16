---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_datasource_urls"
description: |-
  Use this data source to get the list of DataArts Security datasource urls within HuaweiCloud.
---

# huaweicloud_dataarts_security_datasource_urls

Use this data source to get the list of DataArts Security datasource urls within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_security_datasource_urls" "test" {
  workspace_id = var.workspace_id
}
```

### Filter by datasource type

```hcl
variable "workspace_id" {}
variable "cluster_id" {}
variable "datasource_type" {}

data "huaweicloud_dataarts_security_datasource_urls" "test" {
  workspace_id     = var.workspace_id
  cluster_id       = var.cluster_id
  datasource_type  = var.datasource_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the datasource urls are located.  
If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the datasource urls belong.

* `cluster_id` - (Optional, String) Specifies the ID of the cluster.

* `datasource_type` - (Optional, String) Specifies the type of the datasource.

* `parent_permission_set_id` - (Optional, String) Specifies the ID of the parent permission set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `urls` - The list of datasource urls that matched filter parameters.  
  The [urls](#dataarts_security_datasource_urls_attr) structure is documented below.

<a name="dataarts_security_datasource_urls_attr"></a>
The `urls` block supports:

* `name` - The name of the url path.

* `contains` - Whether the parent permission set contains this permission.
