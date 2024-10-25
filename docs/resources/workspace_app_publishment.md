---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_publishment"
description: |-
  Manages a Workspace APP pulishing resource within HuaweiCloud.
---

# huaweicloud_workspace_app_publishment

Manages a Workspace APP pulishing resource within HuaweiCloud.

-> 1. Before using this resource, ensure that the `type` parameter of the `huaweicloud_workspace_app_group` resource
   must be **COMMON_APP** and `server_group_id` parameter must be set.
   <br>2. Deleting this resource will unpublish the APP.

## Example Usage

```hcl
variable "app_group_id" {}
variable "app_name" {}
variable "execute_path" {}

resource "huaweicloud_workspace_app_publishment" "test" {
  app_group_id = var.app_group_id
  name         = var.app_name
  type         = 3
  execute_path = var.execute_path
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `app_group_id` - (Required, String, ForceNew) Specifies the APP group ID to which the application belongs.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the application.  
  The name valid length is limited from `1` to `64` and cannot be all spaces.
  The name must be unique.

* `type` - (Required, Int, ForceNew) Specifies the type of the application.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **2**: Private image APP.
  + **3**: Custom APP.

* `execute_path` - (Required, String) Specifies the location where the application file is installed.
  e.g. `C:\Program Files\Internet Explorer\iexplore.exe`.

* `sandbox_enable` - (Optional, Bool) Specifies whether to run in sandbox mode, defaults to `false`.  
  If you want to set `true`, please ensure that the application sandbox software has been installed on the associated server
  group instance. Otherwise, the application cannot be started.
  
* `version` - (Optional, String) Specifies the version of the application.  
  If the `sandbox_enable` is set to `true`, this parameter value is the version of the sandboxed application.

* `publisher` - (Optional, String, ForceNew) Specifies the publisher of the application.
  Changing this creates a new resource.  
  If the `sandbox_enable` is set to `true`, this parameter value is the publisher of the sandboxed application.

* `work_path` - (Optional, String) Specifies the publisher of the application, e.g. `C:\Program Files\Internet Explorer`.

* `command_param` - (Optional, String) Specifies the command line parameter used to start the application.  
  If the `sandbox_enable` is set to `true`, the path of the APP to be started must be enclosed in
  double quotation marks (""), e.g. `/box:DefaultBox "C:\Program Files\Internet Explorer\iexplore.exe"`.

* `description` - (Optional, String) Specifies the description of the application.

* `icon_index` - (Optional, Int, ForceNew) Specifies the icon index of the application.
  Changing this creates a new resource.

* `icon_path` - (Optional, String, ForceNew) Specifies the path where the application icon is located.
  Changing this creates a new resource.

* `status` - (Optional, String) Specifies the current status of the application, defaults to **NORMAL**.
  The valid values are as follows:
  + **NORMAL**
  + **FORBIDDEN**

* `source_image_ids` - (Optional, List, ForceNew) Specifies the list of image IDs corresponding to the server instance
  to which the application belongs.  
  The maximum length is `20`.
  This parameter is required and available only when the `type` is `2`.  

  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `published_at` - The release time of the application, in RFC3339 format.

## Import

The resource can be imported using `app_group_id` and `name`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_workspace_app_publishment.test <app_group_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `source_image_ids`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_workspace_app_publishment" "test" {
  ...

  lifecycle {
    ignore_changes = [
      source_image_ids,
    ]
  }
}
```
