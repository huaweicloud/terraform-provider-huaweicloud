---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_environment"
description: |-
  Manages a CodeArts deploy environment resource within HuaweiCloud.
---

# huaweicloud_codearts_deploy_environment

Manages a CodeArts deploy environment resource within HuaweiCloud.

## Example Usage

```hcl
variable "project_id" {}
variable "application_id" {}
variable "name" {}
variable "gourp_id" {}
variable "host_id" {}

resource "huaweicloud_codearts_deploy_environment" "test" {
  project_id     = var.project_id
  application_id = var.application_id
  name           = var.name
  deploy_type    = 0
  os_type        = "linux"
  description    = "demo"

  hosts {
    group_id = var.gourp_id
    host_id  = var.host_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, ForceNew) Specifies the project ID for CodeArts service.
  Changing this creates a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the application ID.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the environment name.

* `deploy_type` - (Required, Int, ForceNew) Specifies the deployment type.
  Valid value are as follows:
  + **0**: Host.
  + **1**: Kubernetes.

  Changing this creates a new resource.

* `os_type` - (Required, String, ForceNew) Specifies the operating system.
  **Windows** or **Linux**, which must be the same as that of the host cluster.
  Changing this creates a new resource.

* `description` - (Optional, String) Specifies the description.

* `hosts` - (Optional, List) Specifies the target hosts list.
  The [hosts](#block--hosts) structure is documented below.

  -> If you import a target host bound to a proxy host, the proxy host will be imported to the environment automatically.
  A proxy host is deleted, when its last target host is deleted from the environment.

<a name="block--hosts"></a>
The `hosts` block supports:

* `group_id` - (Required, String) Specifies the cluster group ID.

* `host_id` - (Required, String) Specifies the host ID to be imported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the create time.

* `created_by` - Indicates the creator information.
  The [created_by](#attrblock--created_by) structure is documented below.

* `hosts` - Indicates the hosts list.
  The [hosts](#attrblock--hosts) structure is documented below.

* `instance_count` - Indicates the number of host instances in the environment.

* `permission` - Indicates the user permission.
  The [permission](#attrblock--permission) structure is documented below.

* `permission_matrix` - Indicates the permission matrix.
  The [permission_matrix](#attrblock--permission_matrix) structure is documented below.

* `proxies` - Indicates the proxy hosts list.
  The [proxies](#attrblock--proxies) structure is documented below.

<a name="attrblock--created_by"></a>
The `created_by` block supports:

* `user_id` - Indicates the user ID.

* `user_name` - Indicates the user name.

<a name="attrblock--hosts"></a>
The `hosts` block supports:

* `host_name` - Indicates the host name.

* `ip_address` - Indicates the IP address.

* `connection_status` - Indicates the connection status.

<a name="attrblock--permission"></a>
The `permission` block supports:

* `can_delete` - Indicates whether the user has the permission to delete environments.

* `can_deploy` - Indicates whether the user has the deploy permission.

* `can_edit` - Indicates whether the user has the permission to edit environments.

* `can_manage` - Indicates whether the user has the permission to edit the environment permission matrix.

* `can_view` - Indicates whether the user has the view environment.

<a name="attrblock--permission_matrix"></a>
The `permission_matrix` block supports:

* `can_delete` - Indicates whether the role has the permission to delete environments.

* `can_deploy` - Indicates whether the role has the deploy permission.

* `can_edit` - Indicates whether the role has the permission to edit environments.

* `can_manage` - Indicates whether the role has the permission to edit the environment permission matrix.

* `can_view` - Indicates whether the role has the view environment.

* `permission_id` - Indicates the permission ID.

* `role_id` - Indicates the role ID.

* `role_name` - Indicates the role name.

* `role_type` - Indicates the role type.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.

<a name="attrblock--proxies"></a>
The `proxies` block supports:

* `connection_status` - Indicates the connection status.

* `group_id` - Indicates the cluster group ID.

* `host_id` - Indicates the host ID.

* `host_name` - Indicates the host name.

* `ip_address` - Indicates the IP address.

## Import

The environment can be imported using `project_id`, `application_id`, and `id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_codearts_deploy_environment.test <project_id>/<application_id>/<id>
```
