---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_uniagent_batch_install"
description: |-
  Use this resource to batch install AOM UniAgent within HuaweiCloud.
---

# huaweicloud_aom_uniagent_batch_install

Use this resource to batch install AOM UniAgent within HuaweiCloud.

-> This resource is only a one-time action resource for batch installing AOM UniAgents. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from
   the tfstate file.

## Example Usage

```hcl
variable "proxy_region_id" {}
variable "installer_agent_id" {}
variable "version" {}
variable "agent_import_param_list" {
  type = list(object({
    account  = string
    password = string
    inner_ip = string
    port     = int
    os_type  = string
    agent_id = string
    vpc_id   = string
  }))
}
variable "install_version" {}

resource "huaweicloud_aom_uniagent_batch_install" "test" {
  proxy_region_id      = var.proxy_region_id
  installer_agent_id   = var.installer_agent_id
  version              = var.version
  public_net_flag      = false
  icagent_install_flag = true

  dynamic "agent_import_param_list" {
    for_each = var.agent_import_param_list

    content {
      account  = agent_import_param_list.value.account
      password = agent_import_param_list.value.password
      inner_ip = agent_import_param_list.value.inner_ip
      port     = agent_import_param_list.value.port
      os_type  = agent_import_param_list.value.os_type
      agent_id = agent_import_param_list.value.agent_id
      vpc_id   = agent_import_param_list.value.vpc_id
    }
  }

  plugin_install_base_param {
    install_version = var.install_version
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the target machines to be operated are located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `agent_import_param_list` - (Required, List, NonUpdatable) Specifies the list of machine parameters for installing
  UniAgent.  
  Up to `100` machines are supported.  
  The [agent_import_param_list](#aom_agent_import_param_list) structure is documented below.

* `proxy_region_id` - (Required, Int, NonUpdatable) Specifies the proxy region ID.  
  Fill in `0` for direct access, or fill in the actual proxy region ID for proxy access.

* `installer_agent_id` - (Required, String, NonUpdatable) Specifies the agent ID of the installation machine.

* `version` - (Required, String, NonUpdatable) Specifies the version number of UniAgent to be installed.

* `public_net_flag` - (Required, Bool, NonUpdatable) Specifies whether to use public network access.

* `icagent_install_flag` - (Optional, Bool, NonUpdatable) Specifies whether to install ICAgent plugin.
  The valid values are as follows:
  + **true**: the latest version of ICAgent plugin will be installed by default.
  + **false**: no ICAgent plugin will be installed.

* `plugin_install_base_param` - (Optional, List, NonUpdatable) Specifies the basic information for plugin
  installation.  
  The [plugin_install_base_param](#aom_plugin_install_base_param) structure is documented below.

<a name="aom_agent_import_param_list"></a>
The `agent_import_param_list` block supports:

* `password` - (Required, String) Specifies the login password of the target machine.

* `inner_ip` - (Required, String) Specifies the IP address of the target machine.

* `port` - (Required, Int) Specifies the login port of the target machine, default is `22`.

* `account` - (Required, String) Specifies the SSH account of the target machine.

* `os_type` - (Required, String) Specifies the operating system type of the target machine.
  The valid values are as follows:
  + **LINUX**
  + **WINDOWS**

* `agent_id` - (Optional, String) Specifies the unique value of the agent.

  -> Required if the corresponding machine is already installed this agent and you want to re-import.

* `vpc_id` - (Optional, String) Specifies the VPC ID of the target machine.

* `coc_cmdb_id` - (Optional, String) Specifies the external unique identifier for COC usage.

<a name="aom_plugin_install_base_param"></a>
The `plugin_install_base_param` block supports:

* `install_version` - (Optional, String) Specifies the specified ICAgent version to install.
  + When `icagent_install_flag` is set to **true**:
    - If `plugin_install_base_param` is empty, the latest version of ICAgent plugin will be installed by default.
    - If `plugin_install_base_param` is specified with a version, that version will be installed.
  + When `icagent_install_flag` is set to **false**:
    - No ICAgent plugin will be installed regardless of this parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the batch action resource.
