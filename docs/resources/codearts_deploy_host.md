---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_host"
description: ""
---

# huaweicloud_codearts_deploy_host

Manages a CodeArts deploy host resource within HuaweiCloud.

-> The servers need to be prepared first. Please refer to the following documents to complete your server configuration:
[Preparations](https://support.huaweicloud.com/intl/en-us/usermanual-deployman/deployman_hlp_0018.html),
[Security Configuration](https://support.huaweicloud.com/intl/en-us/usermanual-deployman/deployman_hlp_1103.html),
[Target Host Configuration](https://support.huaweicloud.com/intl/en-us/usermanual-deployman/deployman_hlp_1101.html) and
[Proxy Host Configuring](https://support.huaweicloud.com/intl/en-us/usermanual-deployman/deployman_hlp_1102.html).

## Example Usage

### Creating a proxy host

```hcl
variable "group_id" {}
variable "group_os_type" {}
variable "ip_address" {}
variable "port" {}
variable "username" {}
variable "password" {}

resource "huaweicloud_codearts_deploy_host" "test" {
  group_id   = var.group_id
  ip_address = var.ip_address
  port       = var.port
  username   = var.username
  password   = var.password
  os_type    = var.group_os_type
  name       = "test_proxy_host"
  as_proxy   = true
}
```

### Creating a target host with proxy access mode

```hcl
variable "group_id" {}
variable "group_os_type" {}
variable "ip_address" {}
variable "port" {}
variable "username" {}
variable "password" {}
variable "proxy_host_id" {}

resource "huaweicloud_codearts_deploy_host" "test" {
  group_id      = var.group_id
  ip_address    = var.ip_address
  port          = var.port
  username      = var.username
  password      = var.password
  os_type       = var.group_os_type
  proxy_host_id = var.proxy_host_id
  name          = "test_hostname"
}
```

### Creating a target host without proxy access mode

```hcl
variable "group_id" {}
variable "group_os_type" {}
variable "ip_address" {}
variable "port" {}
variable "username" {}
variable "password" {}

resource "huaweicloud_codearts_deploy_host" "test" {
  group_id   = var.group_id
  ip_address = var.ip_address
  port       = var.port
  username   = var.username
  password   = var.password
  os_type    = var.group_os_type
  name       = "test_hostname"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the CodeArts deploy group ID.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the host name. The name consists of 3 to 128 characters, including letters,
  digits, chinese characters or `-_.` symbols.

* `ip_address` - (Required, String) Specifies the IP address of your server.
  + When creating a proxy host, only public IPv4 addresses are supported.
  + When creating a target host with proxy access mode, both public and private IPv4 addresses are supported.
  + When creating a target host without proxy access mode, only public IPv4 addresses are supported.

* `port` - (Required, Int) Specifies the SSH port of your server. The value ranges from `1` to `65,535`.

* `os_type` - (Required, String, ForceNew) Specifies the operating system. Valid values are **windows** and **linux**.
  The value must be consistent with the CodeArts deploy group.

  Changing this parameter will create a new resource.

* `username` - (Required, String) Specifies the username of your server.

* `password` - (Optional, String) Specifies the password of your server.

* `private_key` - (Optional, String) Specifies the private key of your server.

-> The parameter `username`, `password` and `private_key` are used for login authentication. At least one of
`private_key` and `password` must be set. And the field `private_key` and `password` can not be set together.

* `as_proxy` - (Optional, Bool, ForceNew) Specifies whether the host is an agent host. Defaults to **false**.

  Changing this parameter will create a new resource.

* `proxy_host_id` - (Optional, String) Specifies the proxy host ID. A proxy host ID can be assigned to multiple target
  hosts. This field is required only when creating a target host with proxy access mode.

* `install_icagent` - (Optional, Bool) Specifies whether to enable Application Operations Management (AOM) for free
  to provide metric monitoring, log query and alarm functions. Defaults to **false**.

* `sync` - (Optional, Bool) Specifies whether to synchronize the password of the current host to the hosts with the
  same IP address, username and port number in other group in the same project. Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The create time.

* `updated_at` - The update time.

* `lastest_connection_at` - The last connection time.

* `connection_status` - The connection status. Valid values are **success**, **failed** and **pending**.

  -> The host can be used only when `connection_status` is **success**. If the field is **failed**, please refer to the
  following documents to check your servers:
  [The hosts FAQs](https://support.huaweicloud.com/intl/en-us/deployman_faq/deployman_faq_0000.html) and
  [The environment FAQs](https://support.huaweicloud.com/intl/en-us/deployman_faq/deployman_faq_00001.html)

* `permission` - The host permission detail.
  The [permission](#DeployHost_permission) structure is documented below.

<a name="DeployHost_permission"></a>
The `permission` block supports:

* `can_view` - Indicates whether the user has the view permission.

* `can_edit` - Indicates whether the user has the edit permission.

* `can_delete` - Indicates whether the user has the deletion permission.

* `can_add_host` - Indicates whether the user has the permission to add hosts.

* `can_copy` - Indicates whether the user has the permission to copy hosts.

## Import

The CodeArts deploy host resource can be imported using `group_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_codearts_deploy_host.test <group_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password`, `private_key`, `install_icagent`
and `sync`. It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_deploy_host" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      password,
      private_key,
      install_icagent,
      sync,
    ]
  }
}
```
