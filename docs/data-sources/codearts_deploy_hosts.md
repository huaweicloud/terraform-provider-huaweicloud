---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_hosts"
description: |-
  Use this data source to get the list of CodeArts deploy hosts.
---

# huaweicloud_codearts_deploy_hosts

Use this data source to get the list of CodeArts deploy hosts.

## Example Usage

```hcl
variable "group_id" {}

data "huaweicloud_codearts_deploy_hosts" "test" {
  group_id = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `group_id` - (Required, String) Specifies the group ID.

* `as_proxy` - (Optional, String) Specifies whether the host is proxy or not.
  Valid values are **true** and **false**.

* `environment_id` - (Optional, String) Specifies the environment ID.

* `name` - (Optional, String) Specifies the name of host.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `hosts` - Indicates the host list.
  The [hosts](#attrblock--hosts) structure is documented below.

<a name="attrblock--hosts"></a>
The `hosts` block supports:

* `id` - Indicates the host ID.

* `name` - Indicates the host name.

* `as_proxy` - Indicates whether the host is an agent host.

* `connection_status` - Indicates the connection status.

* `env_count` - Indicates the number of environments.

* `import_status` - Indicates the import status.

* `ip_address` - Indicates the IP address.

* `lastest_connection_at` - Indicates the last connection time.

* `os_type` - Indicates the operating system.

* `owner_id` - Indicates the owner ID.

* `owner_name` - Indicates the owner name.

* `created_at` - Indicates the create time.

* `permission` - Indicates the permission.
  The [permission](#attrblock--hosts--permission) structure is documented below.

* `port` - Indicates the SSH port.

* `proxy_host_id` - Indicates the agent ID.

* `trusted_type` - Indicates the trusted type.
  + **0** indicates password authentication.
  + **1** indicates key authentication.

* `username` - Indicates the username.

<a name="attrblock--hosts--permission"></a>
The `permission` block supports:

* `can_add_host` - Indicates whether the user has the permission to add hosts.

* `can_copy` - Indicates whether the user has the permission to copy hosts.

* `can_delete` - Indicates whether the user has the deletion permission.

* `can_edit` - Indicates whether the user has the edit permission.

* `can_view` - Indicates whether the user has the view permission.
