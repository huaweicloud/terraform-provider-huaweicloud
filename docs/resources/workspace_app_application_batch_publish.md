---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_application_batch_publish"
description: |-
  Use this resource to batch publish applications of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_application_batch_publish

Use this resource to batch publish applications of the Workspace APP within HuaweiCloud.

-> This resource is a one-time action resource used to batch publish applications. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "application_group_id" {}
variable "applications" {
  type = list(object({
    name              = string
    execute_path      = string
    source_type       = number
    sandbox_enable    = optional(bool)
    version           = optional(string)
    publisher         = optional(string)
    work_path         = optional(string)
    command_param     = optional(string)
    description       = optional(string)
    icon_path         = optional(string)
    icon_index        = optional(number)
    source_image_ids  = optional(list(string))
    is_pre_boot       = optional(bool)
    app_extended_info = optional(map(string))
  }))
}

resource "huaweicloud_workspace_app_application_batch_publish" "test" {
  app_group_id = var.application_group_id

  dynamic "applications" {
    for_each = var.applications
    content {
      name              = applications.value["name"]
      execute_path      = applications.value["execute_path"]
      source_type       = applications.value["source_type"]
      sandbox_enable    = applications.value["sandbox_enable"]
      version           = applications.value["version"]
      publisher         = applications.value["publisher"]
      work_path         = applications.value["work_path"]
      command_param     = applications.value["command_param"]
      description       = applications.value["description"]
      icon_path         = applications.value["icon_path"]
      icon_index        = applications.value["icon_index"]
      source_image_ids  = applications.value["source_image_ids"]
      is_pre_boot       = applications.value["is_pre_boot"]
      app_extended_info = applications.value["app_extended_info"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the applications to be published are located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `app_group_id` - (Required, String, NonUpdatable) Specifies the ID of the application group.

* `applications` - (Required, List, NonUpdatable) Specifies the list of applications to be published.  
  The [applications](#app_application_batch_publish_applications) structure is documented below.

<a name="app_application_batch_publish_applications"></a>
The `applications` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the name of the application.  
  The application name allows visible characters or spaces, but cannot be all spaces.  
  The length range is `1` to `64` characters.

* `execute_path` - (Required, String, NonUpdatable) Specifies the execution path of the application.

* `source_type` - (Required, Int, NonUpdatable) Specifies the type of the application.  
  The valid values are as follows:
  + **2**: Private image APP.
  + **3**: Custom APP.

* `version` - (Optional, String, NonUpdatable) Specifies the version of the application.  

* `command_param` - (Optional, String, NonUpdatable) Specifies the command line parameters used to start
  the application.  
  If the `sandbox_enable` is set to **true**, the path of the application to be started must be enclosed in
  double quotation marks (""), e.g. `/box:DefaultBox "C:\Program Files\Internet Explorer\iexplore.exe"`.

* `work_path` - (Optional, String, NonUpdatable) Specifies the working directory of the application.

* `icon_path` - (Optional, String, NonUpdatable) Specifies the path where the application icon is located.

* `icon_index` - (Optional, Int, NonUpdatable) Specifies the icon index of the application.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the application.

* `publisher` - (Optional, String, NonUpdatable) Specifies the publisher of the application.  
    If the `sandbox_enable` is set to **true**, this parameter value is the publisher of the sandboxed application.

* `source_image_ids` - (Optional, List, NonUpdatable) Specifies the list of image IDs to which the application
  belongs.  
  The maximum length is `20`.  
  This parameter is required and available only when the `source_type` is `2`.

* `sandbox_enable` - (Optional, Bool, NonUpdatable) Specifies whether to run in sandbox mode.
  Defaults to **false**.

* `is_pre_boot` - (Optional, Bool, NonUpdatable) Specifies whether to enable application pre-boot.  
  Defaults to **false**.

* `app_extended_info` - (Optional, Map, NonUpdatable) Specifies the extended information of the custom application.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `applications` - The list of published applications.  
  The [applications](#app_application_batch_publish_applications_attr) structure is documented below.

<a name="app_application_batch_publish_applications_attr"></a>
The `applications` block supports:

* `id` - The ID of the published application.
