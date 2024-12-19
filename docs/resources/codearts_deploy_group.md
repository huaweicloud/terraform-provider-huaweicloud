---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_group"
description: ""
---

# huaweicloud_codearts_deploy_group

Manages a CodeArts deploy group resource within HuaweiCloud.

## Example Usage

### Using proxy access mode

```hcl
variable "project_id" {}

resource "huaweicloud_codearts_deploy_group" "test" {
  project_id  = var.project_id
  name        = "test_group"
  os_type     = "linux"
  description = "test description"
}
```

### Without using proxy access mode

```hcl
variable "project_id" {}

resource "huaweicloud_codearts_deploy_group" "test" {
  project_id    = var.project_id
  name          = "test_group"
  os_type       = "linux"
  description   = "test description"
  is_proxy_mode = 0
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the group name. The name consists of 3 to 128 characters, including letters,
  digits, chinese characters or `-_.` symbols.

* `project_id` - (Required, String, ForceNew) Specifies the project ID for CodeArts service.

  Changing this parameter will create a new resource.

* `os_type` - (Required, String, ForceNew) Specifies the operating system. Valid values are **windows** and **linux**.

  Changing this parameter will create a new resource.

* `resource_pool_id` - (Optional, String) Specifies the resource pool ID. A resource pool is a collection of physical
  environments that execute deployment commands when deploying software packages.
  If not specified, the resource pool hosted by HuaweiCloud will be used.
  If you want to use your own servers as resource pools, please fill your own resource pool ID.

  -> Please refer to the following documents to create your own resource pool:
  [Creating an Agent Pool](https://support.huaweicloud.com/intl/en-us/usermanual-devcloud/devcloud_01_0016.html) and
  [Creating an Agent](https://support.huaweicloud.com/intl/en-us/usermanual-devcloud/devcloud_01_0017.html).

* `description` - (Optional, String) Specifies the description.

* `is_proxy_mode` - (Optional, Int, ForceNew) Specifies whether the host is in agent access mode.
  Changing this parameter will create a new resource. Valid values are as follows:
  + **1**: Using proxy access mode.
  + **0**: Without using proxy access mode.

  Defaults to 1.

  -> The scenes to use agent access mode: If the target host cannot connect to the public network, you need to select a
  host bound with an EIP as the proxy host to connect CodeArts to the target host.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (group ID).

* `created_at` - The create time.

* `updated_at` - The update time.

* `created_by` - The creator information.
  The [object](#DeployGroup_user) structure is documented below.

* `permission` - The group permission detail.
  The [permission](#DeployGroup_permission) structure is documented below.

* `permission_matrix` - The group permission matrix detail.
  The [permission_matrix](#DeployGroup_permission_matrix) structure is documented below.

<a name="DeployGroup_user"></a>
The `object` block supports:

* `user_id` - The user ID.

* `user_name` - The user name.

<a name="DeployGroup_permission"></a>
The `permission` block supports:

* `can_view` - Indicates whether the user has the view permission.

* `can_edit` - Indicates whether the user has the edit permission.

* `can_delete` - Indicates whether the user has the deletion permission.

* `can_add_host` - Indicates whether the user has the permission to add hosts.

* `can_manage` - Indicates whether the user has the management permission.

* `can_copy` - Indicates whether the user has the permission to copy.

<a name="DeployGroup_permission_matrix"></a>
The `permission_matrix` block supports:

* `role_id` - Indicates the role ID.

* `role_name` - Indicates the role name.

* `role_type` - Indicates the role type.

* `can_view` - Indicates whether the role has the view permission.

* `can_edit` - Indicates whether the role has the edit permission.

* `can_delete` - Indicates whether the role has the deletion permission.

* `can_add_host` - Indicates whether the role has the permission to add hosts.

* `can_manage` - Indicates whether the role has the management permission.

* `can_copy` - Indicates whether the role has the permission to copy.

* `created_at` - The permission create time.

* `updated_at` - The permission update time.

## Import

The CodeArts deploy group resource can be imported using the `project_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_codearts_deploy_group.test <project_id>/<id>
```
