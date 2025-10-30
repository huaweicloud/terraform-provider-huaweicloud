---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_dashboards_folder"
description:  |-
  Manages an AOM dashboards folder resource within HuaweiCloud.
---

# huaweicloud_aom_dashboards_folder

Manages an AOM dashboards folder resource within HuaweiCloud.

## Example Usage

```hcl
variable "folder_title" {}

resource "huaweicloud_aom_dashboards_folder" "test" {
  folder_title = var.folder_title 
  delete_all   = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `folder_title` - (Required, String) Specifies the dashboards folder title.

* `delete_all` - (Optional, Bool) Specifies whether to delete the dashboards when deleting folder. Defaults to **false**.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the folder belongs.
  Defaults to **0**. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

* `dashboard_ids` - Indicates the dashboard IDs under the folder.

* `is_template` - Indicates whether the folder is default.

* `created_by` - Indicates the creator of the folder.

## Import

The AOM dashboards folder resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_aom_dashboards_folder.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from
the API response. The missing attributes include: `delete_all`.
It is generally recommended running `terraform plan` after importing a dashboards folder.
You can then decide if changes should be applied to the dashboards folder, or the resource definition
should be updated to align with the dashboards folder. Also you can ignore changes as below.

```hcl
resource "huaweicloud_aom_dashboards_folder" "test" {
  ...

  lifecycle {
    ignore_changes = [
      delete_all,
    ]
  }
}
```
