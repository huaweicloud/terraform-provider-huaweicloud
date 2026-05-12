---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_dataservice_catalogs"
description: |-
  Use this data source to query DataArts DataService catalogs within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_catalogs

Use this data source to query DataArts DataService catalogs within HuaweiCloud.

## Example Usage

### Query the catalogs under the root path of a specified workspace

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_dataservice_catalogs" "test" {
  workspace_id = var.workspace_id
}
```

### Query the catalogs under a specified workspace and filter by catalog ID

```hcl
variable "workspace_id" {}
variable "sub_catalog_id" {}

data "huaweicloud_dataarts_dataservice_catalogs" "test" {
  workspace_id = var.workspace_id
  catalog_id   = var.sub_catalog_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the catalogs are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the catalogs belong.

* `catalog_id` - (Optional, String) Specifies the ID of the catalog.  
  Defaults to `0`, which represents the root directory.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `catalogs` - The list of catalogs that matched filter parameters.  
  The [catalogs](#dataarts_dataservice_catalogs) structure is documented below.

<a name="dataarts_dataservice_catalogs"></a>
The `catalogs` block supports:

* `id` - The ID of the catalog, in UUID format.

* `parent_id` - The ID of the parent catalog for the current catalog.

* `name` - The name of the catalog.

* `description` - The description of the catalog.

* `api_catalog_type` - The type of the API catalog.

* `created_at` - The creation time of the catalog, in RFC3339 format.

* `create_user` - The creator of the catalog.

* `updated_at` - The last update time of the catalog, in RFC3339 format.

* `update_user` - The last updater of the catalog.
