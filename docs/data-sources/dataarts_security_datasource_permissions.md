---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_datasource_permissions"
description: |-
  Use this data source to get the list of DataArts Security datasource configurable permissions within HuaweiCloud.
---

# huaweicloud_dataarts_security_datasource_permissions

Use this data source to get the list of DataArts Security datasource configurable permissions within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_security_datasource_permissions" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the datasource permissions are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the datasource permissions belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datasource_permissions` - The list of datasource configurable permissions.  
  The [datasource_permissions](#dataarts_security_datasource_permissions_attr) structure is documented below.

<a name="dataarts_security_datasource_permissions_attr"></a>
The `datasource_permissions` block supports:

* `datasource_type` - The datasource type.

* `permission_types` - The list of datasource operation permission types.

* `permission_actions` - The list of supported datasource permission actions.
