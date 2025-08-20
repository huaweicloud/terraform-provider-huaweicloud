---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_warehouse_application"
description: |-
  Manages an application resource of Workspace APP warehouse within HuaweiCloud.
---

# huaweicloud_workspace_app_warehouse_application

Manages an application resource of Workspace APP warehouse within HuaweiCloud.

## Example Usage

```hcl
var "application_name" {}
var "version" {}
var "version_name" {}
var "file_store_obs_path" {}

resource "huaweicloud_workspace_app_warehouse_application" "test" {
  name            = var.application_name
  category        = "OTHER"
  os_type         = "Windows"
  version         = var.version
  version_name    = var.version_name
  file_store_path = var.file_store_obs_path
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the application.
  The valid length is limited from `1` to `64`, only Chinese and English characters, digits, underscores (_) and
  hyphens (-) are allowed.

* `category` - (Required, String) Specifies the category of the application.  
  The vaild values are as follows:
  + **GAME**
  + **SECURE_STORAGE**
  + **MULTIMEDIA_AND_CODING**
  + **PROJECT_MANAGEMENT**
  + **PRODUCTIVITY_AND_COLLABORATION**
  + **GRAPHIC_DESIGN**
  + **OTHER**

* `os_type` - (Required, String) Specifies the operating system type of the application.
  The valid values are as follows:
  + **Windows**
  + **Linux**
  + **Other**

* `version` - (Required, String) Specifies the version of the application.
  The valid length is limited from `1` to `64` and cannot contain spaces.

* `version_name` - (Required, String) Specifies the version name of the application.
  The valid length is limited from `1` to `64` and cannot contain spaces.

* `file_store_path` - (Required, String, ForceNew) Specifies the storage path of the OBS bucket where the application
  is located. Changing this creates a new resource.

  -> 1.The OBS bucket name where the uploaded file is located must consist of `wks-app` and the project ID, connected by
     a hyphen (-). e.g. `wks-app-0970dd7a1300f5672ff2c003c60ae115`.<br>2.The path is the relative path of the file in the
     OBS bucket. For example, the relative path of the `https:/xxx.xxx.com/file/workspace_app.exe` file is
     `file/workspace_app.exe`.<br>3.Deleting the resource to which this `file_store_path` belongs, the file in the OBS
     storage path will also be deleted synchronously.
  
* `description` - (Optional, String) Specifies the description of the application.

* `icon` - (Optional, String) Specifies the icon of the application.  
  The valid value ranges from `0` to `10,977`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also application ID.

* `record_id` - The record ID of the application.

## Import

The resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_app_warehouse_application.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `icon`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_workspace_app_warehouse_application" "test" {
  ...

  lifecycle {
    ignore_changes = [
      icon,
    ]
  }
}
```
