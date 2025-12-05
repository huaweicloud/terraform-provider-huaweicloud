---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_setting_docker_plugin_install_script"
description: |-
  Use this data source to get docker plug-in installation script.
---

# huaweicloud_hss_setting_docker_plugin_install_script

Use this data source to get docker plug-in installation script.

## Example Usage

```hcl
variable "operate_type" {}

data "huaweicloud_hss_setting_docker_plugin_install_script" "test" {
  plugin       = "opa-docker-authz"
  operate_type = var.operate_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `plugin` - (Required, String) Specifies the plug-in type.
  The valid value is **opa-docker-authz**.

* `operate_type` - (Required, String) Specifies the operation type.
  The valid values are as follows:
  + **install**
  + **upgrade**
  + **uninstall**

* `outside_host` - (Optional, Bool) Specifies whether a server is a non-Huawei Cloud server.
  The valid values are as follows:
  + **true**
  + **false** (Default value)

* `batch_install` - (Optional, Bool) Specifies whether to install in batches.
  The valid values are as follows:
  + **true** (Default value)
  + **false**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The docker plug-in script information.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `package_type` - The installation package type.

* `cmd` - The command cmd.

* `package_download_url` - The package download address.
