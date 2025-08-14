---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application"
description: |-
  Using this resource to manage Workspace application within Huaweicloud.
---

# huaweicloud_workspace_application

Using this resource to manage Workspace application within Huaweicloud.

## Example Usage

### Create an application with OBS bucket store

```hcl
variable "application_name" {}
variable "application_version" {}
variable "application_description" {}
variable "authorization_type" {}
variable "install_type" {}
variable "support_os" {}
variable "catalog_id" {}
variable "bucket_name" {}
variable "bucket_file_path" {}

resource "huaweicloud_workspace_application" "test" {
  name               = var.application_name
  version            = var.application_version
  description        = var.application_description
  authorization_type = var.authorization_type
  install_type       = var.install_type
  support_os         = var.support_os
  catalog_id         = var.catalog_id
  reserve_obs_file   = true

  application_file_store {
    store_type = "OBS"

    bucket_store {
      bucket_name      = var.bucket_name
      bucket_file_path = var.bucket_file_path
    }
  }
}
```

### Create an application with external file link

```hcl
variable "application_name" {}
variable "application_version" {}
variable "application_description" {}
variable "authorization_type" {}
variable "install_type" {}
variable "support_os" {}
variable "catalog_id" {}
variable "file_link" {}

resource "huaweicloud_workspace_application" "test" {
  name               = var.application_name
  version            = var.application_version
  description        = var.application_description
  authorization_type = var.authorization_type
  install_type       = var.install_type
  support_os         = var.support_os
  catalog_id         = var.catalog_id

  application_file_store {
    store_type = "LINK"
    file_link  = var.file_link
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the application.  
  The value can contain `1` to `128` characters.
  Can not consist solely of spaces and the following special characters are not allowed: ``:^;|`~{}[]<>``.

* `version` - (Required, String) Specifies the version of the application.  
  The value can contain `1` to `128` characters.
  Only letters, digits, hyphens(-), underscores(_) and dots(.) are allowed in the version number.

* `description` - (Required, String) Specifies the description of the application.

* `authorization_type` - (Required, String) Specifies the authorization type of the application.  
  The valid values are as follows:
  + **ALL_USER**
  + **ASSIGN_USER**

* `application_file_store` - (Required, List) Specifies the file store configuration of the application.  
  The [application_file_store](#workspace_application_file_store) structure is documented below.

* `install_type` - (Required, String) Specifies the installation type of the application.  
  The valid values are as follows:
  + **QUIET_INSTALL**: Silent installation mode, suitable for automated deployment.
  + **UNZIP_INSTALL**: Extract and install mode, for applications that need to be extracted first.
  + **GUI_INSTALL**: Graphical installation mode, requires user interaction during installation.

* `support_os` - (Required, String) Specifies the supported operating system of the application.  
  The valid values are as follows:
  + **Linux**
  + **Windows**
  + **Other**

* `catalog_id` - (Required, String) Specifies the catalog ID of the application.

* `application_icon_url` - (Optional, String) Specifies the icon URL of the application.

* `install_command` - (Optional, String) Specifies the installation command of the application.

* `uninstall_command` - (Optional, String) Specifies the uninstallation command of the application.

* `install_info` - (Optional, String) Specifies the installation information of the application.

* `reserve_obs_file` - (Optional, Bool) Specifies whether to delete the installation package in the OBS bucket.  
  Required if the value of parameter `store_type` is **OBS**.

<a name="workspace_application_file_store"></a>
The `application_file_store` block supports:

* `store_type` - (Required, String) Specifies the store type of the application file.  
  The valid values are as follows:
  + **OBS**: Object Storage Service bucket store.
  + **LINK**: External file link.

* `bucket_store` - (Optional, List) Specifies the OBS bucket store configuration.  
  The [bucket_store](#workspace_application_bucket_store) structure is documented below.
  Required if the value of parameter `store_type` is **OBS**.

* `file_link` - (Optional, String) Specifies the external file link.  
  Required if the value of parameter `store_type` is **LINK**.

<a name="workspace_application_bucket_store"></a>
The `bucket_store` block supports:

* `bucket_name` - (Required, String) Specifies the name of the OBS bucket.

* `bucket_file_path` - (Required, String) Specifies the file path in the OBS bucket.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the application.

* `application_source` - The source of the application.

* `create_time` - The creation time of the application, in UTC format.

* `catalog` - The catalog name of the application.

## Import

Application can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_application.test <id>
```
