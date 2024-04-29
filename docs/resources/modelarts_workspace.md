---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_workspace"
description: ""
---

# huaweicloud_modelarts_workspace

Manages a Modelarts workspace resource within HuaweiCloud.  

## Example Usage

### Create a public workspace

```hcl
resource "huaweicloud_modelarts_workspace" "test" {
  name        = "demo"
  description = "This is a demo"
}
```

### Create an internal workspace

```hcl
variable "user_id" {}
resource "huaweicloud_modelarts_workspace" "test" {
  name        = "demo"
  description = "This is a demo"
  auth_type   = "internal"
  grants {
    user_id = var.user_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Workspace name, which consists of 4 to 64 characters.  
  Only chinese and english letters, digits, hyphens (-), and underscores (_) are allowed.
  **default** is the name of the default workspace reserved by the system.

* `description` - (Optional, String) The description of the workspace.  

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the workspace.  
  Value 0 indicates the default enterprise project.

  Changing this parameter will create a new resource.

* `auth_type` - (Optional, String) Authorization type.  
  Value options are as follows:
    + **public**: public access within the tenant.
    + **private**: Only the creator and main account can access.
    + **internal**: Creator, main account, and specified IAM sub-accounts can access.

  Defaults to **public**.

* `grants` - (Optional, List) List of authorized users.  
  It is mandatory when **auth_type** is **internal**.
  The [grants](#ModelartsWorkspace_Grants) structure is documented below.

<a name="ModelartsWorkspace_Grants"></a>
The `grants` block supports:

* `user_id` - (Optional, String) IAM user ID.  
  User ID and username specify at least one. If both are specified, User ID is preferred.

* `user_name` - (Optional, String) IAM username.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Workspace status.  
  Valid values are **CREATE_FAILED**, **NORMALL**, **DELETING** and **DELETE_FAILED**.

* `status_info` - Status details.  
  If the deletion fails, you can check the reason through this field.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Default is 10 minutes.

## Import

The Modelarts workspace can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_workspace.test 0ce123456a00f2591fabc00385ff1234
```
