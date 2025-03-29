---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_uniagent"
description: |-
  Manages an AOM UniAgent resource within HuaweiCloud.
---

# huaweicloud_aom_uniagent

Manages an AOM UniAgent resource within HuaweiCloud.

-> Destroying resource does not uninstall the UniAgent.

## Example Usage

### Install UniAgent for a host

```hcl
variable "installer_agent_id" {}
variable "inner_ip" {}
variable "account" {}
variable "password" {}

resource "huaweicloud_aom_uniagent" "test" {
  installer_agent_id = var.installer_agent_id
  version            = "1.1.6"
  public_net_flag    = false
  proxy_region_id    = 0
  inner_ip           = var.inner_ip
  port               = 22
  account            = var.account
  password           = var.password
  os_type            = "LINUX"
}
```

### Reinstall UniAgent for a host

```hcl
variable "installer_agent_id" {}
variable "inner_ip" {}
variable "account" {}
variable "password" {}
variable "agent_id" {}

resource "huaweicloud_aom_uniagent" "test" {
  installer_agent_id = var.installer_agent_id
  version            = "1.1.6"
  public_net_flag    = false
  proxy_region_id    = 0
  agent_id           = var.agent_id
  inner_ip           = var.inner_ip
  port               = 22
  account            = var.account
  password           = var.password
  os_type            = "LINUX"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `installer_agent_id` - (Required, String, NonUpdatable) Specifies the installer agent ID.

* `version` - (Required, String) Specifies the UniAgent version to be installed.

  -> When updating the `version`, `agent_id` must be specified.

* `public_net_flag` - (Required, Bool, NonUpdatable) Specifies the whether to use public network access.

* `proxy_region_id` - (Required, Int, NonUpdatable) Specifies the proxy region ID.
  + Specifies it as **0** to use the direct access.
  + Specifies it as specific proxy region ID to use the proxy access.

* `inner_ip` - (Required, String, NonUpdatable) Specifies the IP of the host where the UniAgent will be installed.

* `port` - (Required, String, NonUpdatable) Specifies the login port of the host where the UniAgent will be installed.

* `account` - (Required, String, NonUpdatable) Specifies the login account of the host where the UniAgent will be
  installed.

* `password` - (Required, String, NonUpdatable) Specifies the login password of the host where the UniAgent will be
  installed.

* `os_type` - (Required, String, NonUpdatable) Specifies the OS type of the host where the UniAgent will be installed.

* `agent_id` - (Optional, String) Specifies the agent ID of the host where the UniAgent will be installed.

  -> If you are reinstalling the UniAgent or updating the UniAgent version, `agent_id` must specified.

* `vpc_id` - (Optional, String, NonUpdatable) Specifies the VPC ID of the host where the UniAgent will be installed.

* `coc_cmdb_id` - (Optional, String, NonUpdatable) Specifies the COC CMDB ID of the host where the UniAgent will be
  installed.

* `icagent_install_flag` - (Optional, Bool, NonUpdatable) Specifies whether to install ICAgent. Defaults to **false**.

* `icagent_install_version` - (Optional, String, NonUpdatable) Specifies the ICAgent version to be installed.
  If it's not specified and `icagent_install_flag` is **true**, install the latest version by default.

* `access_key` - (Optional, String, NonUpdatable) Specifies the access key of the IAM account where the host ICAgent is
  not installed.

* `secret_key` - (Optional, String, NonUpdatable) Specifies the secret key of the IAM account where the host ICAgent is
  not installed.

  -> AK/SK is required only when installing the older version of ICAgent.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
