---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_warehouse_applications"
description: |-
  Use this data source to get warehouse application list of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_warehouse_applications

Use this data source to get warehouse application list of the Workspace APP within HuaweiCloud.

## Example Usage

```hcl
variable "application_id" {}

data "huaweicloud_workspace_app_warehouse_applications" "test" {
  app_id = var.application_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the warehouse applications are located.  
  If omitted, the provider-level region will be used.

* `app_id` - (Optional, String) Specifies the ID of the application.

* `name` - (Optional, String) Specifies the name of the application.  
  Fuzzy matching is supported.

* `category` - (Optional, String) Specifies the category of the application.  
  The valid values are as follows:
  + **GAME**
  + **SECURE_STORAGE**
  + **MULTIMEDIA_AND_CODING**
  + **PROJECT_MANAGEMENT**
  + **PRODUCTIVITY_AND_COLLABORATION**
  + **GRAPHIC_DESIGN**
  + **OTHER**

* `verify_status` - (Optional, String) Specifies the verification status of the application.  
  The valid values are as follows:
  + **VERIFIED** - Verification passed.
  + **VERIFY_FAILED** - Verification failed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `applications` - All applications that match the filter parameters.  
  The [applications](#data_app_warehouse_applications) structure is documented below.

<a name="data_app_warehouse_applications"></a>
The `applications` block supports:

* `id` - The record ID of the application.

* `app_id` - The ID of the application.

* `name` - The name of the application.

* `category` - The category of the application.

* `os_type` - The operating system type of the application.
  + **Windows**
  + **Linux**
  + **Other**

* `version` - The version of the application.

* `version_name` - The version name of the application.

* `file_store_path` - The storage path of the application file.

* `app_file_size` - The size of the application file.

* `description` - The description of the application.

* `verify_status` - The verification status of the application.

* `icon` - The base64 encoded application icon.

* `created_at` - The creation time of the application, in RFC3339 format.

* `updated_at` - The latest update time of the application, in RFC3339 format.
