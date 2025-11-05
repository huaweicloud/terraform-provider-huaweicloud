---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_publishable_applications"
description: |-
  Use this data source to get the list of the publishable applications under specified APP group within HuaweiCloud.
---

# huaweicloud_workspace_app_publishable_applications

Use this data source to get the list of the publishable applications under specified APP group within HuaweiCloud.

## Example Usage

```hcl
variable "app_group_id" {}

data "huaweicloud_workspace_app_availability_apps" "test" {
  app_group_id = var.app_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `app_group_id` - (Required, String) Specifies the ID of the application group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `group_images` - The list of image IDs under the server group.

* `applications` - The list of the publishable applications.

  The [applications](#app_publishable_applications_attr) structure is documented below.

<a name="app_publishable_applications_attr"></a>
The `applications` block supports:

* `name` - The name of the the application.

* `execute_path` - The execution path where the application is located.

* `version` - The version of the the application.

* `publisher` - The publisher of the the application.

* `work_path` - The work path of the the application.

* `command_param` - The command line arguments used to start the application.

* `description` - The description of the the application.

* `publishable` - Whether the application can be published.

* `icon_index` - The index of the application icon.

* `icon_path` - The path where the application icon is located.

* `source_image_ids` - The list of image IDs to which the application belongs.
