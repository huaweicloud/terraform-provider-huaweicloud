---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_apps"
description: |-
  Use this data source to get the list of Data Service applications within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_apps

Use this data source to get the list of Data Service applications within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "app_name_to_be_queried" {}

data "huaweicloud_dataarts_dataservice_apps" "test" {
  workspace_id = var.workspace_id
  name         = var.app_name_to_be_queried
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the applications are located.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the applications belong.

* `dlm_type` - (Optional, String) Specifies the type of DLM engine.  
  The valid values are as follows:
  + **SHARED**: Shared data service.
  + **EXCLUSIVE**: The exclusive data service.

  Defaults to **SHARED**.

* `name` - (Optional, String) Specifies the name of the applications to be fuzzy queried.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `apps` - All applications that match the filter parameters.  
  The [apps](#dataservice_apps_elem) structure is documented below.

<a name="dataservice_apps_elem"></a>
The `apps` block supports:

* `id` - The ID of the application, in UUID format.

* `name` - The name of the application.

* `description` - The description of the application.

* `app_key` - The key of the application.

* `app_secret` - The secret of the application.

* `created_at` - The creation time of the application, in RFC3339 format.

* `updated_at` - The latest update time of the application, in RFC3339 format.

* `create_user` - The name of the application creator.

* `update_user` - The name of the application updater.

* `app_type` - The type of the application.
