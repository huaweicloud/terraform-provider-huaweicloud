---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_applications"
description: |-
  Use this data source to get the list of the Workspace applications within HuaweiCloud.
---

# huaweicloud_workspace_applications

Use this data source to get the list of the Workspace applications within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_workspace_applications" "test" {}
```

### Filter applications by name

```hcl
variable "application_name" {}

data "huaweicloud_workspace_applications" "test" {
  name = var.application_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the applications are located.  
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the application name to be queried and supports fuzzy matching.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `applications` - The list of applications that match the filter parameters.  
  The [applications](#workspace_application_attr) structure is documented below.

<a name="workspace_application_attr"></a>
The `applications` block supports:

* `id` - The ID of the application.

* `name` - The name of the application.

* `description` - The description of the application.

* `version` - The version of the application.

* `authorization_type` - The authorization type of the application.
  + **ALL_USER**
  + **ASSIGN_USER**

* `application_file_store` - The file store configuration of the application.  
  The [application_file_store](#workspace_application_file_store) structure is documented below.

* `application_icon_url` - The icon URL of the application.

* `install_type` - The installation type of the application.
  + **QUIET_INSTALL**
  + **UNZIP_INSTALL**
  + **GUI_INSTALL**

* `install_command` - The installation command of the application.

* `uninstall_command` - The uninstallation command of the application.

* `support_os` - The supported operating system of the application.
  + **Linux**
  + **Windows**
  + **Other**

* `status` - The status of the application.
  + **NORMAL**
  + **FORBIDDEN**

* `application_source` - The source of the application.
  + **CUSTOM**
  + **SYSTEM**
  + **MARKET**

* `create_time` - The creation time of the application, in UTC format.

* `catalog_id` - The catalog ID of the application.

* `catalog` - The catalog name of the application.

* `install_info` - The installation information of the application.

<a name="workspace_application_file_store"></a>
The `application_file_store` block supports:

* `store_type` - The store type of the application file.
  + **OBS**
  + **LINK**

* `bucket_store` - The OBS bucket store configuration.  
  The [bucket_store](#workspace_application_bucket_store) structure is documented below.

* `file_link` - The external file link.

<a name="workspace_application_bucket_store"></a>
The `bucket_store` block supports:

* `bucket_name` - The name of the OBS bucket.

* `bucket_file_path` - The file path in the OBS bucket.
