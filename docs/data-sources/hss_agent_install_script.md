---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_agent_install_script"
description: |-
  Use this data source to get HSS agent install script within HuaweiCloud.
---

# huaweicloud_hss_agent_install_script

Use this data source to get HSS agent install script within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_agent_install_script" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `os_type` - (Required, String) Specifies the OS type. The valid values are **Windows** and **Linux**.

* `os_arch` - (Required, String) Specifies system architecture. The valid values are **x86_64** and **aarch64**.
  + If `os_type` is **Windows**, this field can only set to **x86_64**.

* `outside_host` - (Optional, Bool) Specifies whether it is not HuaweiCloud.

* `batch_install` - (Optional, Bool) Specifies whether to install in bulk.

* `type` - (Optional, String) Specifies type. The valid values are **password** and **ssh_key**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `install_script_list` - The list of installation scripts.  
  The [install_script_list](#install_script_list_struct) structure is documented below.

<a name="install_script_list_struct"></a>
The `install_script_list` block supports:

* `package_type` - The packet type.

* `cmd` - The command.

* `package_download_url` - The package download URL.
