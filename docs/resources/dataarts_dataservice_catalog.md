---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_dataservice_catalog"
description: |-
  Use this resource to manage a shared or exclusive catalog within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_catalog

Use this resource to manage a shared or exclusive catalog within HuaweiCloud.

## Example Usage

### Create a catalog under the ROOT path

```hcl
variable "workspace_id" {}
variable "catalog_name" {}

resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  workspace_id = var.workspace_id
  dlm_type     = "EXCLUSIVE"
  name         = var.catalog_name
}
```

### Create a catalog under the specified catalog

```hcl
variable "workspace_id" {}
variable "parent_catalog_id" {}
variable "catalog_name" {}

resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  workspace_id = var.workspace_id
  parent_id    = var.parent_catalog_id
  dlm_type     = "EXCLUSIVE"
  name         = var.catalog_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the catalog is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the catalog belongs.  
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the catalog.  
  The valid length is limited from `1` to `64`, only Chinese and English characters, digits and hyphens (-) are
  allowed.  
  The name must start with a Chinese or English character or a digit, and the Chinese characters must be in **UTF-8** or
  **Unicode** format.

* `description` - (Optional, String) Specifies the description of the catalog.  
  Maximum of `255` characters are allowed.

* `dlm_type` - (Optional, String, ForceNew) Specifies the type of DLM engine.  
  The valid values are as follows:
  + **SHARED**: Shared data service.
  + **EXCLUSIVE**: The exclusive data service.

  Defaults to **SHARED**. Changing this parameter will create a new resource.

  -> The value of `dlm_type` for all catalogs under this directory must be consistent with this resource.

* `parent_id` - (Optional, String) Specifies the ID of the parent catalog for current catalog.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `path` - The path of current catalog.

* `catalog_total` - The total number of sub-catalogs in the current catalog.

-> The newly added sub-catalogs needs to wait until the next time `terraform refresh` command is executed that the value
   can be refreshed.

* `api_total` - The total number of APIs in the current catalog.

-> The newly added APIs needs to wait until the next time `terraform refresh` command is executed before that value
   can be refreshed.

* `created_at` - The creation time of the catalog.

* `updated_at` - The latest update time of the catalog.

* `create_user` - The creator of the catalog.

* `update_user` - The user who latest updated the catalog.

## Import

The catalog can be imported using `workspace_id`, `dlm_type` and `id` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_dataarts_dataservice_catalog.test <workspace_id>/<dlm_type>/<id>
```

Also, you can omit `dlm_type` and provide just `workspace_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_dataservice_catalog.test <workspace_id>/<id>
```

~> This way only supports importing the catalog of the **SHARED** type, but does not support the catalog imported for
   **EXCLUSIVE** type. If an error is reported, please carefully check the `dlm_type` value to which imported catalog
   you want.

Note that the imported state may not be identical to your resource definition, because the attributes are missing in the
API response. The missing attributes include: `parent_id`.
It is generally recommended running `terraform plan` after importing an resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  ...

  lifecycle {
    ignore_changes = [
      parent_id,
    ]
  }
}
```
